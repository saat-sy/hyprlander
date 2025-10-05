package cli

import (
	"fmt"

	"github.com/saat-sy/hyprlander/pkg/setup"
	"github.com/spf13/cobra"
)

func UpdateCommand() *cobra.Command {
	updateCommand := &cobra.Command{
		Use:   "update",
		Short: "Update gemini api key",
		Long:  "Update the stored gemini api key used for authentication",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Updating Gemini API key...")

			update := setup.NewSetup()

			err := update.Check()
			if err != nil {
				return fmt.Errorf("hyprlander is not initialized. Please run 'hyprlander init' first: %w", err)
			}

			apiKey, err := update.PromptForAPIKey()
			if err != nil {
				return fmt.Errorf("failed to get new API key: %w", err)
			}

			if err := update.Update(apiKey); err != nil {
				return fmt.Errorf("could not update API key: %w", err)
			}

			fmt.Println("Gemini API key updated successfully!")
			return nil
		},
	}

	return updateCommand
}