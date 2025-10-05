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
	context 	context.Context
	chatSession *genai.Chat
	history     []*genai.Content
	maxTurns    int
}

func NewAgent() *Agent {
	history := []*genai.Content{
		genai.NewContentFromText(SystemPrompt, genai.RoleUser),
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

	chat, _ := client.Chats.Create(
		ctx,
		config.GeminiModel,
		tools.Config,
		history,
	)

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

		res, _ := agent.chatSession.SendMessage(agent.context, genai.Part{Text: prompt})

		fmt.Printf("Response: %+v\n", res)
		fmt.Printf("History: %+v\n", agent.chatSession.History(false))

		break
	}
}
