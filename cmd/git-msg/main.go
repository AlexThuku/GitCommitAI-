package main

import (
	"fmt"
	"os"

	"github.com/AlexThuku/GitCommitAI-/internal/ai"
	"github.com/AlexThuku/GitCommitAI-/internal/cli"
	"github.com/AlexThuku/GitCommitAI-/internal/config"
	"github.com/AlexThuku/GitCommitAI-/internal/git"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func main() {
	// Setup logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create AI provider based on configuration
	var provider ai.Provider
	switch cfg.ModelProvider {
	case "openai":
		provider = ai.NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.OpenAIModel)
	case "huggingface":
		provider = ai.NewHuggingFaceProvider(cfg.HuggingFaceToken, cfg.HuggingFaceModel)
	case "local":
		provider = ai.NewLocalProvider(cfg.LocalEndpoint)
	default:
		// Default to Hugging Face if not specified
		provider = ai.NewHuggingFaceProvider(cfg.HuggingFaceToken, cfg.HuggingFaceModel)
	}

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "git-msg",
		Short: "AI-powered Git commit message generator",
		Long:  "Generate meaningful commit messages based on your uncommitted changes using AI",
	}

	// Add generate command
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a commit message",
		Run: func(cmd *cobra.Command, args []string) {
			// Get git diff
			diff, err := git.GetDiff()
			if err != nil {
				slog.Error("Failed to get git diff", "error", err)
				os.Exit(1)
			}

			if diff == "" {
				fmt.Println("No changes detected. Stage your changes first.")
				os.Exit(0)
			}

			fmt.Println("Analyzing changes...")

			// Generate commit message
			message, err := provider.GenerateCommitMessage(diff)
			if err != nil {
				slog.Error("Failed to generate commit message", "error", err)

				// Try fallback if primary provider failed
				var fallbackProvider ai.Provider
				var fallbackName string

				if _, ok := provider.(*ai.HuggingFaceProvider); ok && cfg.LocalEndpoint != "" {
					fallbackName = "local model"
					fallbackProvider = ai.NewLocalProvider(cfg.LocalEndpoint)
				} else if _, ok := provider.(*ai.LocalProvider); ok && cfg.HuggingFaceToken != "" {
					fallbackName = "Hugging Face model"
					fallbackProvider = ai.NewHuggingFaceProvider(cfg.HuggingFaceToken, cfg.HuggingFaceModel)
				} else if _, ok := provider.(*ai.OpenAIProvider); ok {
					if cfg.HuggingFaceToken != "" {
						fallbackName = "Hugging Face model"
						fallbackProvider = ai.NewHuggingFaceProvider(cfg.HuggingFaceToken, cfg.HuggingFaceModel)
					} else if cfg.LocalEndpoint != "" {
						fallbackName = "local model"
						fallbackProvider = ai.NewLocalProvider(cfg.LocalEndpoint)
					}
				}

				if fallbackProvider != nil {
					fmt.Printf("Falling back to %s...\n", fallbackName)
					message, err = fallbackProvider.GenerateCommitMessage(diff)
					if err != nil {
						slog.Error("Fallback failed", "error", err)
						os.Exit(1)
					}
				} else {
					os.Exit(1)
				}
			}

			// Present to user for approval
			approved, finalMessage := cli.PromptForApproval(message)
			if approved {
				if err := git.SetCommitMessage(finalMessage); err != nil {
					slog.Error("Failed to set commit message", "error", err)
					os.Exit(1)
				}
				fmt.Println("Commit message set successfully!")
			} else {
				fmt.Println("Operation cancelled.")
			}
		},
	}

	rootCmd.AddCommand(generateCmd)

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}
