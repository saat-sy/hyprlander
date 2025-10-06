package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/saat-sy/hyprlander/pkg/config"
	"github.com/saat-sy/hyprlander/pkg/core/tools"
	"github.com/saat-sy/hyprlander/pkg/setup"
	"github.com/saat-sy/hyprlander/pkg/ui"
	"google.golang.org/genai"
)

type Agent struct {
	context     context.Context
	chatSession *genai.Chat
	history     []*genai.Content
	maxTurns    int
	ui          ui.UI
}

func NewAgent() *Agent {
	set := setup.NewSetup()
	userUI := ui.New()

	keys, err := set.FetchConfig()
	if err != nil {
		userUI.PrintError(fmt.Errorf("error fetching config: %w", err))
		log.Fatal("Error fetching config:", err)
	}

	apiKey := keys[config.APIKeyName]
	hyprlandDirName := keys[config.HyprlandDirName]

	exists, err := config.DirExists(hyprlandDirName)
	if !exists || err != nil {
		userUI.PrintError(fmt.Errorf("hyprland directory does not exist. Please provide a valid path in the secrets.ini file"))
		log.Fatal("Hyprland directory does not exist. Please provide a valid path in the secrets.ini file.")
	}

	tree, err := config.GetTreeFromDir(hyprlandDirName)
	if err != nil {
		userUI.PrintError(fmt.Errorf("error building directory tree: %w", err))
		log.Fatal("Error building directory tree:", err)
	}

	userUI.Print("System initialized with directory tree")

	history := []*genai.Content{
		genai.NewContentFromText(GetSystemPrompt(tree), genai.RoleUser),
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		userUI.PrintError(err)
		log.Fatal(err)
	}

	tools := tools.NewConfigForTools()

	chat, err := client.Chats.Create(
		ctx,
		config.GeminiModel,
		tools.Config,
		history,
	)
	if err != nil {
		userUI.PrintError(fmt.Errorf("failed to create chat session: %w", err))
		log.Fatalf("Failed to create chat session: %v", err)
	}

	return &Agent{
		context:     ctx,
		chatSession: chat,
		history:     history,
		maxTurns:    10,
		ui:          userUI,
	}
}

func (agent *Agent) InvokeAgent(prompt string) {
	agent.ui.Print(fmt.Sprintf("Initial Request: %s", prompt))

	currentPrompt := prompt
	var pendingFunctionResponse *genai.FunctionResponse

	for turn := 1; turn <= agent.maxTurns; turn++ {
		agent.ui.Print(fmt.Sprintf("----- Turn %d -----", turn))

		var parts []genai.Part

		if pendingFunctionResponse != nil {
			parts = append(parts, genai.Part{FunctionResponse: pendingFunctionResponse})
			pendingFunctionResponse = nil
		}

		if currentPrompt != "" {
			parts = append(parts, genai.Part{Text: currentPrompt})
		}

		res, err := agent.chatSession.SendMessage(agent.context, parts...)
		if err != nil {
			agent.ui.PrintError(fmt.Errorf("error sending message: %w", err))
			continue
		}

		if len(res.Candidates) == 0 {
			agent.ui.Print("No response from the model. Trying again...")
			continue
		}

		candidate := res.Candidates[0]
		if len(candidate.Content.Parts) == 0 {
			agent.ui.Print("No content parts in response. Trying again...")
			continue
		}

		if textPart := candidate.Content.Parts[0].Text; textPart != "" {
			agent.ui.PrintAgent(textPart)
			userInput, err := agent.userInteraction()
			if err != nil {
				agent.ui.PrintError(fmt.Errorf("error during user interaction: %w", err))
				continue
			}

			if userInput != "" {
				currentPrompt = userInput
			} else {
				agent.ui.Print("No user input provided. Ending interaction.")
				return
			}
		} else if funcCall := candidate.Content.Parts[0].FunctionCall; funcCall != nil {
			agent.ui.PrintTool(funcCall.Name, funcCall.Args)

			confirm, err := agent.confirmExecution()
			if err != nil {
				agent.ui.PrintError(fmt.Errorf("error during confirmation: %w", err))
				continue
			}

			if !confirm {
				agent.ui.Print("Function execution cancelled by user.")
				return
			}

			output, err := agent.executeFunctionCall(funcCall)
			if err != nil {
				agent.ui.PrintError(fmt.Errorf("error executing function call: %w", err))
				currentPrompt = fmt.Sprintf("The function call failed with error: %v. Please provide an alternative solution.", err)
				continue
			}

			agent.ui.PrintSuccess(fmt.Sprintf("Function Output: %s", output))

			pendingFunctionResponse = &genai.FunctionResponse{
				Name:     funcCall.Name,
				Response: map[string]interface{}{"result": output},
			}
		}
	}
}

func (agent *Agent) confirmExecution() (bool, error) {
	return agent.ui.Confirm("Do you want to proceed?")
}

func (agent *Agent) userInteraction() (string, error) {
	agent.ui.Print("User Interaction Required:")
	return agent.ui.Input("Please provide any necessary suggestion or leave blank: ")
}

func (agent *Agent) executeFunctionCall(funcCall *genai.FunctionCall) (string, error) {
	switch funcCall.Name {
	case "readFile":
		path, ok := funcCall.Args["path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid path parameter for readFile")
		}
		content, err := tools.ReadFile(path)
		if err != nil {
			return "", err
		}
		return content, nil

	case "writeFile":
		path, ok := funcCall.Args["path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid path parameter for writeFile")
		}
		content, ok := funcCall.Args["content"].(string)
		if !ok {
			return "", fmt.Errorf("invalid content parameter for writeFile")
		}
		err := tools.WriteFile(path, content)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Successfully wrote to file: %s", path), nil

	case "shellExecute":
		command, ok := funcCall.Args["command"].(string)
		if !ok {
			return "", fmt.Errorf("invalid command parameter for shellExecute")
		}
		output, err := tools.ShellExecute(command)
		if err != nil {
			return "", err
		}
		return output, nil

	default:
		return "", fmt.Errorf("unknown function: %s", funcCall.Name)
	}
}
