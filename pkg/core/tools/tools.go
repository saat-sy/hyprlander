package tools

import "google.golang.org/genai"

var Tools = []*genai.Tool{
	FileReaderTool,
	FileWriterTool,
	ShellExecutorTool,
}
