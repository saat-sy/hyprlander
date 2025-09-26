package prompt

import (
	"fmt"

	"github.com/spf13/cobra"
)

func PromptCommand() *cobra.Command {
	promptCommand := &cobra.Command{
		Use:   "prompt",
		Short: "Short description for prompt",
		Long:  "Long Description for prompt",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello world from prompt!")
		},
	}

	return promptCommand
}
