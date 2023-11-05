package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type GPTResponse struct {
	ModifiedFunctions []string `json:"modifiedFunctions"`
	Code              string   `json:"code"`
}

func sytemPrompt() string {
	prompt := `
	note: Don't give any single line of explaination, just give me json.
	You must send your response in json format according to bellow schema.
	{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"modifiedFunctions": {
		"description": "An array of modified function names",
		"type": "array",
		"items": {
			"type": "string"
		}
		},
		"code": {
		"description": "The source code after modification or original code if no modifications",
		"type": "string"
		}
	},
	"required": ["modifiedFunctions", "code"]
	}
	`
	return prompt
}

func sendAndGetResponse(content string) (GPTResponse, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: sytemPrompt(),
				},
				{
					// Role conatins 3 types user, assistant and system
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return GPTResponse{}, err
	}

	var responseJSON GPTResponse
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &responseJSON)
	fmt.Println(resp.Choices[0].Message.Content)
	if err != nil {
		return GPTResponse{}, err
	}

	return responseJSON, nil
}
