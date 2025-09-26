package initialize

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InitCommand() *cobra.Command {
	initCommand := &cobra.Command{
		Use:   "init",
		Short: "Short description for Init",
		Long:  "Long Description for init",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello world from init!")
		},
	}

	return initCommand
}
