package git

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetDiff returns the output of git diff for staged changes
func GetDiff() (string, error) {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return "", err
	}

	// Get staged changes with git diff --staged
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// If no staged changes, get unstaged changes
	if out.String() == "" {
		cmd = exec.Command("git", "diff")
		out.Reset()
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}

	return out.String(), nil
}

// SetCommitMessage sets the given commit message for the next commit
// This doesn't actually commit, but prepares the message
func SetCommitMessage(message string) error {
	// Get git root directory
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	gitDir := strings.TrimSpace(out.String())
	commitMsgPath := filepath.Join(gitDir, "COMMIT_EDITMSG")

	// Write message to file safely
	return os.WriteFile(commitMsgPath, []byte(message), 0644)
}

// Commit performs the git commit using the prepared message
func Commit() error {
	cmd := exec.Command("git", "commit", "-F", ".git/COMMIT_EDITMSG")
	return cmd.Run()
}
