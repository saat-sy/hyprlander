package cli

import (
	"fmt"

	"github.com/saat-sy/hyprlander/pkg/core/agent"
	"github.com/spf13/cobra"
)

func PromptCommand() *cobra.Command {
	promptCommand := &cobra.Command{
		Use:   "prompt",
		Short: "Execute prompt-based hyprland configuration changes",
		Long:  "Use natural language prompts to modify hyprland configuration files",
		Run: func(cmd *cobra.Command, args []string) {
			// agent.InvokeAgent(args[0])
		},
	}

	return promptCommand
}
