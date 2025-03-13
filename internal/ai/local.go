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

// LocalProvider implements the Provider interface using a local FastAPI endpoint
type LocalProvider struct {
	endpoint string
	client   *http.Client
}

// LocalRequestPayload defines the payload for the local model API
type LocalRequestPayload struct {
	Diff string `json:"diff"`
}

// LocalResponse defines the response from the local model API
type LocalResponse struct {
	CommitMessage string `json:"commit_message"`
	Error         string `json:"error,omitempty"`
}

// NewLocalProvider creates a new local model provider
func NewLocalProvider(endpoint string) *LocalProvider {
	return &LocalProvider{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateCommitMessage generates a commit message based on the diff
func (p *LocalProvider) GenerateCommitMessage(diff string) (string, error) {
	if diff == "" {
		return "", errors.New("empty diff provided")
	}

	if p.endpoint == "" {
		return "", errors.New("local endpoint URL is not set")
	}

	payload := LocalRequestPayload{
		Diff: diff,
	}

	reqJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", p.endpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("local API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the full response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API response: %w", err)
	}

	var localResp LocalResponse
	if err := json.Unmarshal(body, &localResp); err != nil {
		return "", fmt.Errorf("failed to parse API response: %w (body: %s)", err, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := "Local API error"
		if localResp.Error != "" {
			errMsg = localResp.Error
		}
		return "", errors.New(errMsg)
	}

	if localResp.CommitMessage == "" {
		return "", errors.New("empty response from local API")
	}

	return localResp.CommitMessage, nil
}
