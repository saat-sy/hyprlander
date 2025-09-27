package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/saat-sy/hyprlander/pkg/setup"
	"github.com/spf13/cobra"
)

func promptForAPIKey() (string, error) {
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

func InitCommand() *cobra.Command {
	initCommand := &cobra.Command{
		Use:   "init",
		Short: "Initialize hyprlander configuration",
		Long:  "Initialize hyprlander by setting up configuration directories and API key storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing Hyprlander...")

			init := setup.NewSetup()

			status, _ := init.Check()
			if status {
				fmt.Println("Hyprland already initialized!")
				return nil
			}

			apiKey, err := promptForAPIKey()
			if err != nil {
				return fmt.Errorf("failed to get API key: %w", err)
			}

			if err := init.Run(apiKey); err != nil {
				return fmt.Errorf("could not complete initialization: %w", err)
			}

			fmt.Println("Hyprlander initialized successfully!")
			return nil
		},
	}

	return initCommand
}
