package tools

import (
	"fmt"
	"os"

	"google.golang.org/genai"
)

func WriteFile(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	return nil
}

// TODO: Make this change with respect to the changes required in the file and not overwrite the whole file
var FileWriterTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "writeFile",
			Description: "Writes content to a file at a given path. Creates the file if it does not exist, and overwrites it if it does.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"path": {
						Type:        genai.TypeString,
						Description: "The path of the file to write to.",
					},
					"content": {
						Type:        genai.TypeString,
						Description: "The content to write into the file.",
					},
				},
				Required: []string{"path", "content"},
			},
		},
	},
}
