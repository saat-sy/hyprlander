package cli

import (
	"fmt"

	"github.com/saat-sy/hyprlander/pkg/config"
	"github.com/saat-sy/hyprlander/pkg/setup"
	"github.com/saat-sy/hyprlander/pkg/ui"
	"github.com/spf13/cobra"
)

func InitCommand() *cobra.Command {
	initCommand := &cobra.Command{
		Use:   "init",
		Short: "Initialize hyprlander configuration",
		Long:  "Initialize hyprlander by setting up configuration directories and API key storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			userUI := ui.New()
			userUI.Print("Initializing Hyprlander...")

			init := setup.NewSetup()

			err := init.Check()
			if err == nil {
				userUI.PrintSuccess("Hyprland already initialized!")
				return nil
			}

			apiKey, err := init.Prompt("Please enter your Gemini API key: ")
			if err != nil {
				return fmt.Errorf("failed to get API key: %w", err)
			}

			hyprlandDir, err := init.Prompt("Please enter the path to your Hyprland configuration directory (e.g., /home/user/.config/hypr): ")
			if err != nil {
				return fmt.Errorf("failed to get Hyprland config directory: %w", err)
			}

			values := map[string]string{
				config.APIKeyName:      apiKey,
				config.HyprlandDirName: hyprlandDir,
			}
			if err := init.Run(values); err != nil {
				return fmt.Errorf("could not complete initialization: %w", err)
			}

			userUI.PrintSuccess("Hyprlander initialized successfully!")
			return nil
		},
	}

	return initCommand
}
