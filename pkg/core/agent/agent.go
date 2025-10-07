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

const defaultMaxTurns = 10

func NewAgent() *Agent {
	agent := &Agent{
		context:  context.Background(),
		maxTurns: defaultMaxTurns,
		ui:       ui.New(),
	}

	if err := agent.initialize(); err != nil {
		log.Fatal("Failed to initialize agent:", err)
	}

	return agent
}

func (a *Agent) initialize() error {
	keys, hyprlandDir, err := a.setupConfiguration()
	if err != nil {
		return fmt.Errorf("setup configuration failed: %w", err)
	}

	tree, err := a.validateAndGetDirectoryTree(hyprlandDir)
	if err != nil {
		return fmt.Errorf("directory validation failed: %w", err)
	}

	if err := a.createChatSession(keys[config.APIKeyName], tree); err != nil {
		return fmt.Errorf("chat session creation failed: %w", err)
	}

	return nil
}

func (a *Agent) setupConfiguration() (map[string]string, string, error) {
	set := setup.NewSetup()
	keys, err := set.FetchConfig()
	if err != nil {
		return nil, "", fmt.Errorf("error fetching config: %w", err)
	}

	hyprlandDir := keys[config.HyprlandDirName]
	if hyprlandDir == "" {
		return nil, "", fmt.Errorf("hyprland directory name not configured")
	}

	return keys, hyprlandDir, nil
}

func (a *Agent) validateAndGetDirectoryTree(hyprlandDir string) ([]string, error) {
	exists, err := config.DirExists(hyprlandDir)
	if err != nil {
		return nil, fmt.Errorf("error checking directory existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("hyprland directory does not exist: %s", hyprlandDir)
	}

	tree, err := config.GetTreeFromDir(hyprlandDir)
	if err != nil {
		return nil, fmt.Errorf("error building directory tree: %w", err)
	}

	return tree, nil
}

func (a *Agent) createChatSession(apiKey string, tree []string) error {
	client, err := genai.NewClient(a.context, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	a.history = []*genai.Content{
		genai.NewContentFromText(GetSystemPrompt(tree), genai.RoleUser),
	}

	tools := tools.NewConfigForTools()
	chat, err := client.Chats.Create(
		a.context,
		config.GeminiModel,
		tools.Config,
		a.history,
	)
	if err != nil {
		return fmt.Errorf("failed to create chat session: %w", err)
	}

	a.chatSession = chat
	return nil
}
