package tools

import (
	"fmt"
	"os"

	"google.golang.org/genai"
)

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}
	return string(content), nil
}

var FileReaderTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "readFile",
			Description: "Reads the entire content of a file given its path and returns it as a string.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"path": {
						Type:        genai.TypeString,
						Description: "The path of the file to read.",
					},
				},
				Required: []string{"path"},
			},
		},
	},
}
