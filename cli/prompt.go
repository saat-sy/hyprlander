package cli

import (
	"fmt"

	"github.com/saat-sy/hyprlander/pkg/config"
	"github.com/saat-sy/hyprlander/pkg/core/agent"
	"github.com/spf13/cobra"
)

func PromptCommand() *cobra.Command {
	promptCommand := &cobra.Command{
		Use:   "prompt",
		Short: "Execute prompt-based hyprland configuration changes",
		Long:  "Use natural language prompts to modify hyprland configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			hyprPath, err := config.GetHyprInstallationPath()
			if err != nil {
				return fmt.Errorf("error getting Hypr installation path :%v. Please install hyprland before you can use hyprlander", err)
			}

			hyprlandPathStructure, err := config.GetTreeFromDir(hyprPath)
			if err != nil {
				return fmt.Errorf("error reading Hypr installation directory: %v", err)
			}

			agent := agent.NewAgent(hyprlandPathStructure)
			if len(args) > 0 {
				prompt := args[0]
				agent.InvokeAgent(prompt)
			}
			return nil
		},
	}

	return promptCommand
}
