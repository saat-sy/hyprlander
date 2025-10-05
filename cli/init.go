package cli

import (
	"fmt"

	"github.com/saat-sy/hyprlander/pkg/setup"
	"github.com/spf13/cobra"
)

func InitCommand() *cobra.Command {
	initCommand := &cobra.Command{
		Use:   "init",
		Short: "Initialize hyprlander configuration",
		Long:  "Initialize hyprlander by setting up configuration directories and API key storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing Hyprlander...")

			init := setup.NewSetup()

			err := init.Check()
			if err == nil {
				fmt.Println("Hyprland already initialized!")
				return nil
			}

			apiKey, err := init.PromptForAPIKey()
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
