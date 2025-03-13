package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openAIEndpoint = "https://api.openai.com/v1/chat/completions"

// OpenAIProvider implements the Provider interface using OpenAI's API
type OpenAIProvider struct {
	apiKey string
	model  string
	client *http.Client
}

// OpenAIRequest represents a request to OpenAI's API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
}

// OpenAIMessage represents a message in the OpenAI API request
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a response from OpenAI's API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateCommitMessage generates a commit message based on the diff
func (p *OpenAIProvider) GenerateCommitMessage(diff string) (string, error) {
	if diff == "" {
		return "", errors.New("empty diff provided")
	}

	if p.apiKey == "" {
		return "", errors.New("OpenAI API key is not set")
	}

	prompt := `You are a helpful assistant that generates git commit messages based on code diffs.
Please analyze the following git diff and generate a clear, concise commit message following the Conventional Commits specification.

The format should be: <type>[optional scope]: <description>

Where <type> is one of:
- feat: A new feature
- fix: A bug fix
- docs: Documentation changes
- style: Changes that don't affect code functionality (formatting, etc.)
- refactor: Code changes that neither fix bugs nor add features
- perf: Performance improvements
- test: Adding or correcting tests
- chore: Changes to build process, dependencies, etc.

The description should be concise but descriptive, written in imperative mood.
Do not include any explanations, just return the commit message.

Git diff:
` + diff

	reqBody := OpenAIRequest{
		Model: p.model,
		Messages: []OpenAIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the full response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API response: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to parse API response: %w (body: %s)", err, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := "OpenAI API error"
		if openAIResp.Error.Message != "" {
			errMsg = openAIResp.Error.Message
		}
		return "", errors.New(errMsg)
	}

	if len(openAIResp.Choices) == 0 {
		return "", errors.New("no response from OpenAI API")
	}

	message := openAIResp.Choices[0].Message.Content
	return message, nil
}
