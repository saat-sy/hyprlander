package cli

import (
	"github.com/saat-sy/hyprlander/pkg/core/agent"
	"github.com/spf13/cobra"
)

func PromptCommand() *cobra.Command {
	promptCommand := &cobra.Command{
		Use:   "prompt",
		Short: "Execute prompt-based hyprland configuration changes",
		Long:  "Use natural language prompts to modify hyprland configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			agent := agent.NewAgent()
			if len(args) > 0 {
				prompt := args[0]
				agent.InvokeAgent(prompt)
			}
			return nil
		},
	}

	return promptCommand
}
