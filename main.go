package main

import (
	"fmt"
	"os"

	"github.com/saat-sy/hyprlander/cli"
)

func main() {
	rootCmd := cli.RootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
