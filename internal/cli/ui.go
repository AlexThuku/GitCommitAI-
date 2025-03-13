package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptForApproval asks the user to approve, edit, or reject the generated commit message
// It returns whether the message was approved and the final message
func PromptForApproval(message string) (bool, string) {
	fmt.Printf("Suggested commit: \"%s\"\n", message)
	fmt.Print("[a]ccept, [e]dit, [r]eject? ")

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "a", "accept":
			return true, message
		case "e", "edit":
			return true, promptForEdit(message)
		case "r", "reject":
			return false, ""
		default:
			fmt.Print("Please enter [a]ccept, [e]dit, or [r]eject: ")
		}
	}
}

// promptForEdit allows the user to edit the message
func promptForEdit(message string) string {
	fmt.Println("Edit your commit message (press Enter when done):")
	fmt.Printf("> %s\n", message) // Show current message as starting point
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return message // Return original on error
	}

	edited := strings.TrimSpace(input)
	if edited == "" {
		return message // Keep original if empty
	}

	return edited
}
