package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		fmt.Println("Input cannot be empty. Please try again.")
	}
}

func (c *Console) Confirm(prompt string) (bool, error) {
	for {
		fmt.Printf("%s (y/n): ", prompt)
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
			fmt.Println("Please enter 'y' or 'n'")
		}
	}
}

func (c *Console) Select(prompt string, options []string) (int, error) {
	fmt.Println(prompt)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	for {
		fmt.Print("Select an option: ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			return -1, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid number")
			continue
		}

		if choice < 1 || choice > len(options) {
			fmt.Printf("Please enter a number between 1 and %d\n", len(options))
			continue
		}

		return choice - 1, nil
	}
}

func (c *Console) Print(message string) {
	fmt.Println(message)
}

func (c *Console) PrintAgent(message string) {
	fmt.Printf("🤖 Agent: %s\n", message)
}

func (c *Console) PrintTool(toolName string, args map[string]interface{}) {
	fmt.Printf("🔧 Tool: %s", toolName)
	if len(args) > 0 {
		fmt.Printf(" with args: %v", args)
	}
	fmt.Println()
}

func (c *Console) PrintError(err error) {
	fmt.Printf("❌ Error: %v\n", err)
}

func (c *Console) PrintSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}
