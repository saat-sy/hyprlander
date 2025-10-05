package cli

import (
	"fmt"
	"strconv"

	"github.com/saat-sy/hyprlander/pkg/setup"
	"github.com/spf13/cobra"
)

func UpdateCommand() *cobra.Command {
	updateCommand := &cobra.Command{
		Use:   "update",
		Short: "Update gemini api key",
		Long:  "Update the stored gemini api key used for authentication",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Fetching current config...")

			update := setup.NewSetup()

			err := update.Check()
			if err != nil {
				return fmt.Errorf("hyprlander is not initialized. Please run 'hyprlander init' first: %w", err)
			}

			currentConfig, err := update.FetchConfig()
			if err != nil {
				return fmt.Errorf("failed to fetch current config: %w", err)
			}

			fmt.Println("What do you want to update?")

			index := 1
			for key := range currentConfig {
				fmt.Printf("%d. %s\n", index, key)
				index++
			}

			inp, err := update.Prompt("Enter the number corresponding to the field you want to update: ")
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			selectedIndex, err := strconv.Atoi(inp)
			if err != nil {
				return fmt.Errorf("invalid input, please enter a number: %w", err)
			}

			if selectedIndex < 1 || selectedIndex >= len(currentConfig) {
				return fmt.Errorf("invalid selection")
			}

			contentMap := make(map[string]string)

			index = 1
			for key, value := range currentConfig {
				if index == selectedIndex {
					newValue, err := update.Prompt(fmt.Sprintf("Enter new value for %s (current: %s): ", key, value))
					if err != nil {
						return fmt.Errorf("failed to read new value: %w", err)
					}
					contentMap[key] = newValue
				} else {
					contentMap[key] = value
				}
				index++
			}

			if err := update.Update(contentMap); err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}

			fmt.Println("Config updated successfully!")
			return nil
		},
	}

	return updateCommand
}
