package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/git-msg/internal/ai"
	"github.com/yourusername/git-msg/internal/cli"
	"github.com/yourusername/git-msg/internal/config"
	"github.com/yourusername/git-msg/internal/git"
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
	if cfg.UseLocalModel {
		provider = ai.NewLocalProvider(cfg.LocalEndpoint)
	} else {
		provider = ai.NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.OpenAIModel)
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

				// Try fallback if OpenAI provider failed
				if _, ok := provider.(*ai.OpenAIProvider); ok && cfg.LocalEndpoint != "" {
					fmt.Println("Falling back to local model...")
					fallbackProvider := ai.NewLocalProvider(cfg.LocalEndpoint)
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
