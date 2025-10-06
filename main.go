package main

import (
	"os"

	"github.com/saat-sy/hyprlander/cli"
	"github.com/saat-sy/hyprlander/pkg/ui"
)

func main() {
	rootCmd := cli.RootCommand()

	if err := rootCmd.Execute(); err != nil {
		userUI := ui.New()
		userUI.PrintError(err)
		os.Exit(1)
	}
}
