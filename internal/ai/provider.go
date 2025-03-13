package ai

// Provider defines the interface for AI-based commit message generation
type Provider interface {
	GenerateCommitMessage(diff string) (string, error)
}
