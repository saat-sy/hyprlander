package agent

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/saat-sy/hyprlander/pkg/config"
	"github.com/saat-sy/hyprlander/pkg/core/tools"
	"github.com/saat-sy/hyprlander/pkg/setup"
	"google.golang.org/genai"
)

type Agent struct {
	context     context.Context
	chatSession *genai.Chat
	history     []*genai.Content
	maxTurns    int
}

func NewAgent() *Agent {
	set := setup.NewSetup()

	keys, err := set.FetchConfig()
	if err != nil {
		log.Fatal("Error fetching config:", err)
	}

	apiKey := keys[config.APIKeyName]
	hyprlandDirName := keys[config.HyprlandDirName]

	exists, err := config.DirExists(hyprlandDirName)
	if !exists || err != nil {
		log.Fatal("Hyprland directory does not exist. Please provide a valid path in the secrets.ini file.")
	}

	tree, err := config.GetTreeFromDir(hyprlandDirName)
	if err != nil {
		log.Fatal("Error building directory tree:", err)
	}

	fmt.Println(GetSystemPrompt(tree))

	history := []*genai.Content{
		genai.NewContentFromText(GetSystemPrompt(tree), genai.RoleUser),
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
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
		log.Fatalf("Failed to create chat session: %v", err)
	}

	return &Agent{
		context:     ctx,
		chatSession: chat,
		history:     history,
		maxTurns:    10,
	}
}

func (agent *Agent) InvokeAgent(prompt string) {
	fmt.Printf("Initial Request: %s\n\n", prompt)

	currentPrompt := prompt
	var pendingFunctionResponse *genai.FunctionResponse

	for turn := 1; turn <= agent.maxTurns; turn++ {
		fmt.Printf("----- Turn %d -----\n", turn)

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
			log.Printf("Error sending message: %v. Trying again...", err)
			continue
		}

		if len(res.Candidates) == 0 {
			fmt.Println("No response from the model. Trying again...")
			continue
		}

		candidate := res.Candidates[0]
		if len(candidate.Content.Parts) == 0 {
			fmt.Println("No content parts in response. Trying again...")
			continue
		}

		if textPart := candidate.Content.Parts[0].Text; textPart != "" {
			fmt.Printf("Model Response: %s\n", textPart)
			userInput, err := agent.userInteraction()
			if err != nil {
				log.Printf("Error during user interaction: %v. Trying again...", err)
				continue
			}

			if userInput != "" {
				currentPrompt = userInput
			} else {
				fmt.Println("No user input provided. Ending interaction.")
				return
			}
		} else if funcCall := candidate.Content.Parts[0].FunctionCall; funcCall != nil {
			fmt.Printf("The agent would like to call a function: %s with args %v\n", funcCall.Name, funcCall.Args)

			confirm, err := agent.confirmExecution()
			if err != nil {
				log.Printf("Error during confirmation: %v. Trying again...", err)
				continue
			}

			if !confirm {
				fmt.Println("Function execution cancelled by user.")
				return
			}

			output, err := agent.executeFunctionCall(funcCall)
			if err != nil {
				log.Printf("Error executing function call: %v. Trying again...", err)
				currentPrompt = fmt.Sprintf("The function call failed with error: %v. Please provide an alternative solution.", err)
				continue
			}

			fmt.Printf("Function Output: %s\n", output)

			pendingFunctionResponse = &genai.FunctionResponse{
				Name:     funcCall.Name,
				Response: map[string]interface{}{"result": output},
			}
		}
	}
}

func (agent *Agent) confirmExecution() (bool, error) {
	for {
		fmt.Print("Do you want to proceed? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed to read data: %w", err)
		}

		data = strings.TrimSpace(data)
		if data == "" {
			fmt.Println("Input cannot be empty. Please enter 'y' or 'n'.")
			continue
		}

		switch data {
		case "y", "Y":
			return true, nil
		case "n", "N":
			return false, nil
		default:
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
}

func (agent *Agent) userInteraction() (string, error) {
	fmt.Println("User Interaction Required:")
	fmt.Println("Please provide any necessary suggestion or leave blank:")
	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read data: %w", err)
	}

	data = strings.TrimSpace(data)
	return data, nil
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
