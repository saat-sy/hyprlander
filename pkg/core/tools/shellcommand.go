package tools

import (
	"fmt"
	"os/exec"
	"strings"

	"google.golang.org/genai"
)

func ShellExecute(command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}
	head := parts[0]
	args := parts[1:]

	cmd := exec.Command(head, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %s, output: %s", err, string(output))
	}
	return string(output), nil
}

var ShellExecutorTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "shellExecute",
			Description: "Executes a shell command and returns the output.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"command": {
						Type:        genai.TypeString,
						Description: "The shell command to execute (e.g., 'ls -l').",
					},
				},
				Required: []string{"command"},
			},
		},
	},
}
