package pkg

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenaiClient struct {
	client *openai.Client
}

func NewOpenaiClient(token string) *OpenaiClient {
	return &OpenaiClient{openai.NewClient(token)}
}

func NewOpenaiClientWithIp(token, ip string) *OpenaiClient {
	config := openai.DefaultConfig(token)
	config.BaseURL = ip
	return &OpenaiClient{openai.NewClientWithConfig(config)}
}

// The `RequestOpenAi` function is a method of the `OpenaiClient` struct. It takes a `message` string
// as input and returns a string and an error.
func (openaiClient *OpenaiClient) RequestOpenAi(message string) (string, error) {
	resp, err := openaiClient.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)
	if err != nil {
		return fmt.Sprintf("ChatCompletion error: %v\n", err), err
	}

	return resp.Choices[0].Message.Content, err
}

// func main() {
// 	client := openai.NewClient("your token")
// 	resp, err := client.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model: openai.GPT4,
// 			Messages: []openai.ChatCompletionMessage{
// 				{
// 					Role:    openai.ChatMessageRoleUser,
// 					Content: "Hello!",
// 				},
// 			},
// 		},
// 	)

// 	if err != nil {
// 		fmt.Printf("ChatCompletion error: %v\n", err)
// 		return
// 	}

// 	fmt.Println(resp.Choices[0].Message.Content)
// }
