package main

import (
	"fmt"

	root "github.com/saat-sy/hyprlander/cmd/root"
)

func main() {
	rootCmd := root.RootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
