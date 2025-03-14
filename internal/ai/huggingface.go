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

const huggingFaceEndpoint = "https://api-inference.huggingface.co/models/"

// HuggingFaceProvider implements the Provider interface using Hugging Face's Inference API
type HuggingFaceProvider struct {
	token   string
	modelID string
	client  *http.Client
}

// HuggingFaceRequest represents a request to Hugging Face's API
type HuggingFaceRequest struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// NewHuggingFaceProvider creates a new Hugging Face provider
func NewHuggingFaceProvider(token, modelID string) *HuggingFaceProvider {
	return &HuggingFaceProvider{
		token:   token,
		modelID: modelID,
		client: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for model inference
		},
	}
}

// GenerateCommitMessage generates a commit message based on the diff
func (p *HuggingFaceProvider) GenerateCommitMessage(diff string) (string, error) {
	if diff == "" {
		return "", errors.New("empty diff provided")
	}

	if p.token == "" {
		return "", errors.New("Hugging Face API token is not set")
	}

	// Create prompt for the model to generate a conventional commit message
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
Only output the commit message, no additional text.

Git diff:
` + diff

	// Prepare request
	reqBody := HuggingFaceRequest{
		Inputs: prompt,
		Parameters: map[string]interface{}{
			"max_new_tokens":   100,
			"temperature":      0.7,
			"top_p":            0.95,
			"do_sample":        true,
			"return_full_text": false,
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create request to Hugging Face API
	endpoint := huggingFaceEndpoint + p.modelID
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.token)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Hugging Face API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle rate limiting
	if resp.StatusCode == 429 {
		return "", errors.New("Hugging Face API rate limit exceeded. Please try again later")
	}

	// Handle other error codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Hugging Face API error (%d): %s", resp.StatusCode, string(body))
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API response: %w", err)
	}

	// Parse response - the structure depends on the model
	var result []string
	if err := json.Unmarshal(body, &result); err != nil {
		// Try alternative response format
		var singleResult string
		if err := json.Unmarshal(body, &singleResult); err != nil {
			var mapResult []map[string]interface{}
			if err := json.Unmarshal(body, &mapResult); err != nil {
				return "", fmt.Errorf("failed to parse API response: %w (body: %s)", err, string(body))
			}

			if len(mapResult) > 0 && mapResult[0]["generated_text"] != nil {
				if text, ok := mapResult[0]["generated_text"].(string); ok {
					return text, nil
				}
			}

			return "", fmt.Errorf("unexpected API response format: %s", string(body))
		}
		return singleResult, nil
	}

	if len(result) == 0 {
		return "", errors.New("empty response from Hugging Face API")
	}

	return result[0], nil
}
