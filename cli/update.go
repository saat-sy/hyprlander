package cli

import (
	"fmt"

	"github.com/saat-sy/hyprlander/pkg/setup"
	"github.com/saat-sy/hyprlander/pkg/ui"
	"github.com/spf13/cobra"
)

func UpdateCommand() *cobra.Command {
	updateCommand := &cobra.Command{
		Use:   "update",
		Short: "Update gemini api key",
		Long:  "Update the stored gemini api key used for authentication",
		RunE: func(cmd *cobra.Command, args []string) error {
			userUI := ui.New()
			userUI.Print("Fetching current config...")

			update := setup.NewSetup()

			err := update.Check()
			if err != nil {
				return fmt.Errorf("hyprlander is not initialized. Please run 'hyprlander init' first: %w", err)
			}

			currentConfig, err := update.FetchConfig()
			if err != nil {
				return fmt.Errorf("failed to fetch current config: %w", err)
			}

			var keys []string
			for key := range currentConfig {
				keys = append(keys, key)
			}

			selectedIndex, err := userUI.Select("What do you want to update?", keys)
			if err != nil {
				return fmt.Errorf("failed to read selection: %w", err)
			}

			selectedKey := keys[selectedIndex]

			contentMap := make(map[string]string)

			for key, value := range currentConfig {
				if key == selectedKey {
					newValue, err := userUI.Input(fmt.Sprintf("Enter new value for %s (current: %s): ", key, value))
					if err != nil {
						return fmt.Errorf("failed to read new value: %w", err)
					}
					contentMap[key] = newValue
				} else {
					contentMap[key] = value
				}
			}

			if err := update.Update(contentMap); err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}

			userUI.PrintSuccess("Config updated successfully!")
			return nil
		},
	}

	return updateCommand
}
