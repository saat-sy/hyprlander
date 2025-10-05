package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/saat-sy/hyprlander/pkg/config"
	"github.com/saat-sy/hyprlander/pkg/core/tools"
	"google.golang.org/genai"
)

type Agent struct {
	context     context.Context
	chatSession *genai.Chat
	history     []*genai.Content
	maxTurns    int
}

func NewAgent(tree []string) *Agent {
	history := []*genai.Content{
		genai.NewContentFromText(GetSystemPrompt(tree), genai.RoleUser),
	}

	ctx := context.Background()
	apiKey, err := config.GetAPIKey()
	if err != nil {
		log.Fatal("Error retrieving API key:", err)
	}
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

	for turn := 1; turn <= agent.maxTurns; turn++ {
		fmt.Printf("----- Turn %d -----\n", turn)

		res, err := agent.chatSession.SendMessage(agent.context, genai.Part{Text: prompt})

		if err != nil {
			log.Printf("Error sending message: %v. Trying again...", err)
			continue
		}

		// to call function: res.Candidates[0].Content.Parts[0].FunctionCall

		if len(res.Candidates) > 0 {
			// Call execute task function
		} else {
			fmt.Println("No response from the model. Trying again...")
			continue
		}
	}
}
