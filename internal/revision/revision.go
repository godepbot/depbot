package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	gitDirectory = ".git"
)

func FindLatestHash(pwd string) (string, error) {
	headData, err := readHead(pwd)
	if err != nil {
		return "", err
	}

	/// HEAD points to a branch that is the current branch.
	if strings.Contains(headData, "refs/") {
		return hashFromBranch(pwd, headData)
	}

	/// HEAD points to a commit in particular.
	if len(headData) > 0 {
		return headData, nil
	}

	return "", fmt.Errorf("Oops! No hash available.")
}

func readHead(pwd string) (string, error) {
	headPath := filepath.Join(pwd, gitDirectory, "HEAD")
	bytes, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("Error is %w", err)
	}

	return strings.TrimSpace(string(bytes)), nil
}

func hashFromBranch(pwd string, headContent string) (string, error) {
	components := strings.Split(headContent, ": ")

	if len(components) < 2 {
		return "", fmt.Errorf("Oops! No hash available.")
	}

	branchPath := components[1]
	fullBranchPath := filepath.Join(pwd, gitDirectory, branchPath)
	bytes, err := os.ReadFile(fullBranchPath)
	if err != nil {
		return "", fmt.Errorf("Oops! Something happened! %v", err)
	}

	return string(bytes), nil
}
