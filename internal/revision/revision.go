package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	git      = ".git"
	head     = "HEAD"
	refsPath = "refs/"
)

func FindLatestHash(pwd string) (string, error) {
	headData, err := readHead(pwd)
	if err != nil {
		fmt.Println("Error is :)", err)
		return "", err
	}

	/// HEAD points to a branch that is the current branch.
	if strings.Contains(headData, refsPath) {
		return hashFromBranch(pwd, headData)
	}

	/// HEAD points to a commit in particular.
	if len(headData) > 0 {
		return headData, nil
	}

	return "", fmt.Errorf("Oops! No hash available.")
}

func readHead(pwd string) (string, error) {
	headPath := filepath.Join(pwd, git, head)
	bytes, err := os.ReadFile(headPath)
	if err != nil {
		fmt.Println("Error is :)", err)
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

func hashFromBranch(pwd string, headContent string) (string, error) {
	components := strings.Split(headContent, ": ")

	expectedComponentes := 2
	if !(len(components) >= expectedComponentes) {
		return "", fmt.Errorf("Oops! No hash available.")
	}

	branchPath := components[1]
	fullBranchPath := filepath.Join(pwd, git, branchPath)
	bytes, err := os.ReadFile(fullBranchPath)
	if err != nil {
		return "", fmt.Errorf("Oops! Something happened! %v", err)
	}

	return string(bytes), nil
}
