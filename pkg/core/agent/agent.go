package agent

import (
	"context"
	"log"

	"github.com/saat-sy/hyprlander/pkg/config"
	"github.com/saat-sy/hyprlander/pkg/core/tools"
	"google.golang.org/genai"
)

type Agent struct {
	chatSession *genai.Chat
	history     []*genai.Content
	maxTurns    int
}

func NewAgent() *Agent {
	history := []*genai.Content{
		genai.NewContentFromText(SystemPrompt, genai.RoleUser),
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
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
		chatSession: chat,
		history:     history,
		maxTurns:    10,
	}
}

func (agent *Agent) InvokeAgent(prompt string) {
	
}
