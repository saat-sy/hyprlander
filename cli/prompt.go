package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func PromptCommand() *cobra.Command {
	promptCommand := &cobra.Command{
		Use:   "prompt",
		Short: "Execute prompt-based hyprland configuration changes",
		Long:  "Use natural language prompts to modify hyprland configuration files",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Prompt functionality coming soon!")
		},
	}

	return promptCommand
}
