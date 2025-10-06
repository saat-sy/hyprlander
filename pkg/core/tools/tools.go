package tools

import (
	"google.golang.org/genai"
)

type Tools struct {
	Config *genai.GenerateContentConfig
}

func NewConfigForTools() *Tools {
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			FileReaderTool,
			FileWriterTool,
			ShellExecutorTool,
		},
	}

	return &Tools{
		Config: config,
	}
}
