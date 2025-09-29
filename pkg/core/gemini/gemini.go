package gemini

import (
	"context"

	"google.golang.org/genai"
)

func GeminiChat(prompt string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey:  "haha",
			Backend: genai.BackendGeminiAPI,
		},
	)

	if err != nil {
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"model",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
