package openai

import (
	"context"
	"net/http"
)

type ChatCompletionOptions struct {
	// ID of the model to use.
	Model Model `json:"model" binding:"required"`
	// The messages to generate chat completions for, in the chat format.
	Messages []ChatMessage `json:"messages" binding:"required"`
	// What sampling temperature to use, between 0 and 2.
	// Higher values like 0.8 will make the output more random, while lower values
	// like 0.2 will make it more focused and deterministic.
	Temperature float32 `json:"temperature,omitempty"`
	// An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p probability mass.
	// So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	TopP float32 `json:"top_p,omitempty"`
	// How many chat completions to generate for each input message.
	N int `json:"n,omitempty"`
	// Up to 4 sequences where the API will stop generating further tokens.
	Stop []string `json:"stop,omitempty"`
	// The maximum number of tokens to generate in the chat completion.
	// The total length of input tokens and generated tokens is limited by the model's context length.
	MaxTokens int `json:"max_tokens,omitempty"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on whether
	// they appear in the text so far, increasing the model's likelihood to talk about new topics.
	PresencePenalty float32 `json:"presence_penalty,omitempty"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing
	// frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim.
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`
}

type ChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatCompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Message      ChatMessage `json:"message"`
		Index        int         `json:"index"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatCompletion given messages, the model will return one or more predicted chat completions.
//
// Docs: https://beta.openai.com/docs/api-reference/chat
func (e *Engine) ChatCompletion(ctx context.Context, opts *ChatCompletionOptions) (*ChatCompletionResponse, error) {
	if err := e.validate.StructCtx(ctx, opts); err != nil {
		return nil, err
	}
	uri := e.apiBaseURL + "/chat/completions"
	if opts.MaxTokens == 0 {
		opts.MaxTokens = defaultMaxTokens
	}
	r, err := marshalJson(opts)
	if err != nil {
		return nil, err
	}
	req, err := e.newReq(ctx, http.MethodPost, uri, "json", r)
	if err != nil {
		return nil, err
	}
	resp, err := e.doReq(req)
	if err != nil {
		return nil, err
	}
	var result ChatCompletionResponse
	if err := unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
