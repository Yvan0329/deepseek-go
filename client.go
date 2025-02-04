package deepseek

import (
	"bufio"
	"context"
	"fmt"

	handlers "github.com/cohesion-org/deepseek-go/handlers"
	utils "github.com/cohesion-org/deepseek-go/utils"
)

type Client struct {
	AuthToken string
	BaseURL   string
}

// NewClient creates a new client with an authentication token and an optional custom baseURL.
// If no baseURL is provided, it defaults to "https://api.deepseek.com/".
func NewClient(AuthToken string, baseURL ...string) *Client {
	url := "https://api.deepseek.com/"
	if len(baseURL) > 0 {
		url = baseURL[0]
	}
	return &Client{
		AuthToken: AuthToken,
		BaseURL:   url,
	}

}

// CreateChatCompletion sends a chat completion request and returns the generated response.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request *ChatCompletionRequest,
) (*handlers.ChatCompletionResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("chat/completions").
		SetBodyFromStruct(request).
		Build(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	resp, err := handlers.HandleSendChatCompletionRequest(req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	updatedResp, err := handlers.HandleChatCompletionResponse(resp)

	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return updatedResp, err
}

// CreateStreamChatCompletion send a chat completion request with stream = true and returns the delta
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request *StreamChatCompletionRequest,
) (ChatCompletionStream, error) {

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("chat/completions").
		SetBodyFromStruct(request).
		BuildStream(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := handlers.HandleSendChatCompletionRequest(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &chatCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}
