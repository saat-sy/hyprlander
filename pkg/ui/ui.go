package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"

	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	BgRed   = "\033[41m"
	BgGreen = "\033[42m"
)

type UI interface {
	Input(prompt string) (string, error)
	InputRequired(prompt string) (string, error)
	Confirm(prompt string) (bool, error)
	Select(prompt string, options []string) (int, error)

	Print(message string)
	PrintAgent(message string)
	PrintTool(toolName string, args map[string]interface{})
	PrintError(err error)
	PrintSuccess(message string)
	PrintWarning(message string)
	PrintTitle(title string)
	PrintSeparator()
}

type Console struct {
	reader *bufio.Reader
}

func New() UI {
	return &Console{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (c *Console) Input(prompt string) (string, error) {
	fmt.Printf("%s%s❯%s ", Cyan, Bold, Reset)
	fmt.Print(prompt)
	input, err := c.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(input), nil
}

func (c *Console) InputRequired(prompt string) (string, error) {
	for {
		input, err := c.Input(prompt)
		if err != nil {
			return "", err
		}
		if input != "" {
			return input, nil
		}
		fmt.Printf("%s%s⚠ Input cannot be empty. Please try again.%s\n", Yellow, Bold, Reset)
	}
}

func (c *Console) Confirm(prompt string) (bool, error) {
	for {
		fmt.Printf("%s%s❯%s %s %s[y/n]%s: ", Cyan, Bold, Reset, prompt, Gray, Reset)
		input, err := c.reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		switch input {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Printf("%s%s⚠ Please enter 'y' or 'n'%s\n", Yellow, Bold, Reset)
		}
	}
}

func (c *Console) Select(prompt string, options []string) (int, error) {
	c.PrintTitle(prompt)
	fmt.Println()

	for i, option := range options {
		fmt.Printf("  %s%s%d)%s %s\n", Cyan, Bold, i+1, Reset, option)
	}
	fmt.Println()

	for {
		fmt.Printf("%s%s❯%s Select an option %s[1-%d]%s: ", Cyan, Bold, Reset, Gray, len(options), Reset)
		input, err := c.reader.ReadString('\n')
		if err != nil {
			return -1, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s%s⚠ Please enter a valid number%s\n", Yellow, Bold, Reset)
			continue
		}

		if choice < 1 || choice > len(options) {
			fmt.Printf("%s%s⚠ Please enter a number between 1 and %d%s\n", Yellow, Bold, len(options), Reset)
			continue
		}

		return choice - 1, nil
	}
}

func (c *Console) Print(message string) {
	fmt.Printf("  %s\n", message)
}

func (c *Console) PrintAgent(message string) {
	fmt.Printf("\n%s%s🤖 Agent:%s\n", Blue, Bold, Reset)
	fmt.Printf("%s\n\n", message)
}

func (c *Console) PrintTool(toolName string, args map[string]interface{}) {
	fmt.Printf("%s%s🔧 Tool:%s %s%s%s", Magenta, Bold, Reset, Cyan, toolName, Reset)
	if len(args) > 0 {
		fmt.Printf("%s%s%v%s", Gray, Dim, args, Reset)
	}
	fmt.Println()
}

func (c *Console) PrintError(err error) {
	fmt.Printf("\n%s%s❌ Error:%s %s\n\n", Red, Bold, Reset, err.Error())
}

func (c *Console) PrintSuccess(message string) {
	fmt.Printf("%s%s✅ %s%s\n", Green, Bold, message, Reset)
}

func (c *Console) PrintWarning(message string) {
	fmt.Printf("%s%s⚠ Warning:%s %s\n", Yellow, Bold, Reset, message)
}

func (c *Console) PrintTitle(title string) {
	fmt.Printf("\n%s%s═══ %s ═══%s\n", Cyan, Bold, title, Reset)
}

func (c *Console) PrintSeparator() {
	fmt.Printf("%s%s────────────────────────────────────────%s\n", Gray, Dim, Reset)
}
