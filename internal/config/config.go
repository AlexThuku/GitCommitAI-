package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// OpenAI settings (keeping for backwards compatibility)
	OpenAIAPIKey string `mapstructure:"openai_api_key"`
	OpenAIModel  string `mapstructure:"openai_model"`

	// Hugging Face settings
	HuggingFaceToken string `mapstructure:"huggingface_token"`
	HuggingFaceModel string `mapstructure:"huggingface_model"`

	// Local model settings
	UseLocalModel bool   `mapstructure:"use_local_model"`
	LocalEndpoint string `mapstructure:"local_endpoint"`

	// General settings
	ModelProvider string `mapstructure:"model_provider"` // "openai", "huggingface", or "local"
}

// Load reads config from file and environment variables
func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("openai_model", "gpt-4o")
	viper.SetDefault("huggingface_model", "mistralai/Mistral-7B-Instruct-v0.2")
	viper.SetDefault("model_provider", "huggingface") // Default to HuggingFace
	viper.SetDefault("use_local_model", false)
	viper.SetDefault("local_endpoint", "http://localhost:8000/generate")

	// Check for config in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".config"))
	}

	// Check for config in current directory
	viper.AddConfigPath(".")

	// Set config name and format
	viper.SetConfigName("git-msg")
	viper.SetConfigType("yaml")

	// Bind environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GIT_MSG")
	viper.BindEnv("openai_api_key")
	viper.BindEnv("huggingface_token")

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, will use defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Validate() error {
	switch c.ModelProvider {
	case "openai":
		if c.OpenAIAPIKey == "" {
			// Check environment variable as backup
			apiKey := os.Getenv("OPENAI_API_KEY")
			if apiKey != "" {
				c.OpenAIAPIKey = apiKey
			} else {
				return errors.New("OpenAI API key is required when using OpenAI provider")
			}
		}
	case "huggingface":
		if c.HuggingFaceToken == "" {
			// Check environment variable as backup
			token := os.Getenv("HUGGINGFACE_TOKEN")
			if token != "" {
				c.HuggingFaceToken = token
			} else {
				return errors.New("Hugging Face token is required when using Hugging Face provider")
			}
		}
	case "local":
		if c.LocalEndpoint == "" {
			return errors.New("local endpoint URL is required when using local model")
		}
	default:
		return fmt.Errorf("invalid model provider: %s", c.ModelProvider)
	}
	return nil
}
