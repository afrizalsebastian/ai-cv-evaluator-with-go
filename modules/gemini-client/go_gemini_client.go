package geminiclient

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type IGeminiClient interface {
	GenerateContent(ctx context.Context, jobTitle, prompt string) (string, error)
}

type geminiClient struct {
	cli   *genai.Client
	model string
}

func NewGeminiAiCLient(ctx context.Context, apiKey, model string) (IGeminiClient, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey, Backend: genai.BackendGeminiAPI})
	if err != nil {
		return nil, err
	}

	return &geminiClient{cli: client, model: model}, nil
}

func (g *geminiClient) GenerateContent(ctx context.Context, jobTitle, prompt string) (string, error) {
	systemInstruction := fmt.Sprintf("You are the head recruiter on company and want to evaluate CV and Project for role %s", jobTitle)
	temp := float32(0.9)
	topP := float32(0.9)
	topK := float32(40.0)
	maxOutputToken := int32(1024)

	config := &genai.GenerateContentConfig{
		TopP:              &topP,
		TopK:              &topK,
		Temperature:       &temp,
		MaxOutputTokens:   maxOutputToken,
		SystemInstruction: genai.NewContentFromText(systemInstruction, genai.RoleModel),
	}

	resp, err := g.cli.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}
