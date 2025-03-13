package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	OpenAIAPIKey  string `mapstructure:"openai_api_key"`
	OpenAIModel   string `mapstructure:"openai_model"`
	UseLocalModel bool   `mapstructure:"use_local_model"`
	LocalEndpoint string `mapstructure:"local_endpoint"`
}

// Load reads config from file and environment variables
func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("openai_model", "gpt-4o")
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
	if !config.UseLocalModel && config.OpenAIAPIKey == "" {
		// Check environment variable as backup
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey != "" {
			config.OpenAIAPIKey = apiKey
		} else {
			return nil, errors.New("OpenAI API key is required when not using local model")
		}
	}

	return &config, nil
}
