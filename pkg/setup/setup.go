package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/saat-sy/hyprlander/pkg/config"
)

type Setup struct{}

func NewSetup() *Setup {
	return &Setup{}
}

func (s *Setup) Run(apiKey string) error {
	configDir, err := config.GetUserHomeDirectory()
	if err != nil {
		return fmt.Errorf("could not determine hyprlander directory: %w", err)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", configDir, err)
	}

	secretFilePath, err := config.GetSecretFilePath()
	if err != nil {
		return fmt.Errorf("could not determine hyprlander folder: %w", err)
	}

	if err := os.WriteFile(secretFilePath, []byte(apiKey), 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}

	return nil
}

func (s *Setup) Check() error {
	secretFilePath, err := config.GetSecretFilePath()
	if err != nil {
		return fmt.Errorf("could not determine hyprlander folder and the secret file: %w", err)
	}

	content, err := os.ReadFile(secretFilePath)
	if err != nil {
		return fmt.Errorf("failed to read secret file: %w", err)
	}
	
	apiKey := string(content)
	if len(apiKey) == 0 {
		return fmt.Errorf("secret file is empty")
	}

	return nil
}

func (s *Setup) Update(apiKey string) error {
	secretFilePath, err := config.GetSecretFilePath()
	if err != nil {
		return fmt.Errorf("could not determine hyprlander folder and the secret file: %w", err)
	}
	
	if err := os.WriteFile(secretFilePath, []byte(apiKey), 0600); err != nil {
		return fmt.Errorf("failed to update secret file: %w", err)
	}

	return nil
}

func (s *Setup) PromptForAPIKey() (string, error) {
	fmt.Print("Please enter your Gemini API key: ")
	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", fmt.Errorf("API key cannot be empty")
	}

	return apiKey, nil
}