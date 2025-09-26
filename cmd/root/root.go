package cmd

import (
	"fmt"

	"github.com/saat-sy/hyprlander/cmd/initialize"
	"github.com/saat-sy/hyprlander/cmd/prompt"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hyprlander",
		Short: "An agent that can modify how hyprland looks!",
		Long:  "Use this package to just give prompts and to directly make changes to the hypr config files",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello world!")
		},
	}

	rootCmd.AddCommand(prompt.PromptCommand())
	rootCmd.AddCommand(initialize.InitCommand())

	return rootCmd
}
