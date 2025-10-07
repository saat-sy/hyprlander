package agent

import (
	"fmt"

	"github.com/saat-sy/hyprlander/pkg/core/tools"
	"google.golang.org/genai"
)

func (a *Agent) executeFunctionCall(funcCall *genai.FunctionCall) (string, error) {
	switch funcCall.Name {
	case "readFile":
		return a.executeReadFile(funcCall.Args)
	case "writeFile":
		return a.executeWriteFile(funcCall.Args)
	case "shellExecute":
		return a.executeShellCommand(funcCall.Args)
	default:
		return "", fmt.Errorf("unknown function: %s", funcCall.Name)
	}
}

func (a *Agent) executeReadFile(args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("invalid path parameter for readFile")
	}

	content, err := tools.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return content, nil
}

func (a *Agent) executeWriteFile(args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("invalid path parameter for writeFile")
	}

	content, ok := args["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content parameter for writeFile")
	}

	if err := tools.WriteFile(path, content); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return fmt.Sprintf("Successfully wrote %d bytes to file: %s", len(content), path), nil
}

func (a *Agent) executeShellCommand(args map[string]interface{}) (string, error) {
	command, ok := args["command"].(string)
	if !ok {
		return "", fmt.Errorf("invalid command parameter for shellExecute")
	}

	output, err := tools.ShellExecute(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute shell command '%s': %w", command, err)
	}

	return output, nil
}
