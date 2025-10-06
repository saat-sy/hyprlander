package setup

import (
	"fmt"
	"os"
	"strings"

	"github.com/saat-sy/hyprlander/pkg/config"
	"github.com/saat-sy/hyprlander/pkg/ui"
)

type Setup struct{
	ui ui.UI
}

func NewSetup() *Setup {
	return &Setup{
		ui: ui.New(),
	}
}

func (s *Setup) Run(values map[string]string) error {
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

	var iniContent strings.Builder
	for key, value := range values {
		iniContent.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}

	if err := os.WriteFile(secretFilePath, []byte(iniContent.String()), 0600); err != nil {
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

func (s *Setup) Update(values map[string]string) error {
	secretFilePath, err := config.GetSecretFilePath()
	if err != nil {
		return fmt.Errorf("could not determine hyprlander folder and the secret file: %w", err)
	}

	var iniContent strings.Builder
	for key, value := range values {
		iniContent.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}

	if err := os.WriteFile(secretFilePath, []byte(iniContent.String()), 0600); err != nil {
		return fmt.Errorf("failed to update secret file: %w", err)
	}

	return nil
}

func (s *Setup) Prompt(message string) (string, error) {
	return s.ui.InputRequired(message)
}

func (s *Setup) FetchConfig() (map[string]string, error) {
	secretFilePath, err := config.GetSecretFilePath()
	if err != nil {
		return nil, fmt.Errorf("could not determine hyprlander folder and the secret file: %w", err)
	}

	content, err := os.ReadFile(secretFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	dataMap := make(map[string]string)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			dataMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return dataMap, nil
}
