package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saat-sy/hyprlander/pkg/config"
)

type Setup struct{}

func NewSetup() *Setup {
	return &Setup{}
}

func (s *Setup) Run(apiKey string) error {
	configDir, err := config.GetHomeDirectory()
	if err != nil {
		return fmt.Errorf("could not determine app directory: %w", err)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", configDir, err)
	}

	secretFilePath := filepath.Join(configDir, config.SecretFileName)
	if err := os.WriteFile(secretFilePath, []byte(apiKey), 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}

	return nil
}

func (s *Setup) Check() (bool, error) {
	configDir, err := config.GetHomeDirectory()
	if err != nil {
		return false, fmt.Errorf("could not determine app directory: %w", err)
	}

	if dirExists, err := config.DirExists(configDir); !dirExists {
		return false, err
	}

	return true, nil
}
