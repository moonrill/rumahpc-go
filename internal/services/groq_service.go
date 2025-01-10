package services

import (
	"encoding/json"
	"strings"

	"github.com/magicx-ai/groq-go/groq"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/types"
)

func GeneratePC() (*types.PC, error) {
	prompt := `Generate a compatible computer build with the following components: processor, motherboard, ram, vga, storage, psu, casing, keyboard, mouse and monitor. Ensure that all components are compatible with each other. Return the response strictly in JSON format with the following keys: processor, motherboard, ram, vga, storage, psu, casing, keyboard, mouse, monitor, and price_est(IDR). Provide no other message, give me the response in json code only`

	req := groq.ChatCompletionRequest{
		Messages: []groq.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:  "gemma2-9b-it",
		Stream: false,
	}

	resp, err := config.GroqClient.CreateChatCompletion(req)
	if err != nil {
		return nil, err
	}

	rawContent := resp.Choices[0].Message.Content
	rawContent = strings.TrimSpace(rawContent)
	rawContent = strings.Trim(rawContent, "```")

	var pc types.PC
	if err := json.Unmarshal([]byte(rawContent), &pc); err != nil {
		return nil, err
	}

	return &pc, nil
}
