# Git-Msg: AI-Powered Commit Message Generator

A CLI tool to generate meaningful Git commit messages using AI.

## Installation

### Using Go

```bash
go install github.com/yourusername/git-msg/cmd/git-msg@latest
```

### Homebrew (macOS)

```bash
brew tap yourusername/git-msg
brew install git-msg
```

### Linux

Download the appropriate .deb or .rpm package from the releases page.

```bash
# For Debian/Ubuntu
sudo dpkg -i git-msg_1.0.0_amd64.deb

# For Red Hat/Fedora
sudo rpm -i git-msg-1.0.0.x86_64.rpm
```

## Configuration

Create a configuration file at `~/.config/git-msg.yaml` or in your project directory:

```yaml
# OpenAI API configuration
openai_api_key: "your-api-key-here" # Or set OPENAI_API_KEY environment variable
openai_model: "gpt-4o"

# Local model configuration
use_local_model: false
local_endpoint: "http://localhost:8000/generate"
```

## Usage

```bash
# Generate a commit message for staged changes
git-msg generate

# For help
git-msg --help
```
```

## Unit Tests

Here's an example test for the git diff parsing functionality:

```go:internal/git/diff_test.go
package git

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDiff(t *testing.T) {
	// Create a temporary directory for a git repo
	tempDir, err := os.MkdirTemp("", "git-msg-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repo
	os.Chdir(tempDir)
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Create a file
	testFile := "test.txt"
	err = os.WriteFile(testFile, []byte("initial content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Add and commit the file
	exec.Command("git", "add", testFile).Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()

	// Modify the file
	err = os.WriteFile(testFile, []byte("modified content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test unstaged changes
	diff, err := GetDiff()
	assert.NoError(t, err)
	assert.Contains(t, diff, "modified content")

	// Stage the changes
	exec.Command("git", "add", testFile).Run()

	// Test staged changes
	diff, err = GetDiff()
	assert.NoError(t, err)
	assert.Contains(t, diff, "modified content")
}
```

## Conclusion

This implementation provides a comprehensive, production-ready Git commit message generator that meets all the requirements. The code follows Go best practices with proper error handling, logging, and modular design. The CLI interface is user-friendly, and the tool supports both OpenAI and local model options with a fallback mechanism.

Key features of this implementation:

1. **Modularity**: The code is organized into logical packages for different components.
2. **Configurability**: Supports configuration via files and environment variables.
3. **Robustness**: Includes error handling and fallback mechanisms.
4. **User Experience**: Provides an interactive interface for reviewing and editing messages.

To build and test the tool, clone the repository, run `go mod tidy` to install dependencies, and then `go build ./cmd/git-msg`. # GitCommitAI-
