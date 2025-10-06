package cli

import (
	"github.com/saat-sy/hyprlander/pkg/ui"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hyprlander",
		Short: "An agent that can modify how hyprland looks!",
		Long:  "Use this package to just give prompts and to directly make changes to the hypr config files",
		Run: func(cmd *cobra.Command, args []string) {
			userUI := ui.New()
			userUI.Print("Welcome to hyprlander! Use 'hyprlander --help' to see available commands.")
		},
	}

	rootCmd.AddCommand(PromptCommand())
	rootCmd.AddCommand(InitCommand())
	rootCmd.AddCommand(UpdateCommand())

	return rootCmd
}
