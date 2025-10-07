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
			Description: "Writes content to a file at a given path. Creates the file if it does not exist, and overwrites it if it does. MUST be used when making configuration changes that require modifying file contents. Always call this tool when the user requests changes that need to be saved to files.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"path": {
						Type:        genai.TypeString,
						Description: "The path of the file to write to.",
					},
					"content": {
						Type:        genai.TypeString,
						Description: "The complete content to write into the file, including both modified and unchanged parts.",
					},
				},
				Required: []string{"path", "content"},
			},
		},
	},
}
