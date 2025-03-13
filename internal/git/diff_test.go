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
