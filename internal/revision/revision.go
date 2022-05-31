package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	git                 = ".git"
	head                = "HEAD"
	refsPath            = "refs/"
	expectedComponentes = 2
)

func FindLatestHash(pwd string) (string, error) {
	path := filepath.Join(pwd, git, head)
	bytes, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Error is :)", err)
		return "", err
	}

	headContent := strings.TrimSpace(string(bytes))

	/// HEAD points to a branch that is the current branch.
	if strings.Contains(headContent, refsPath) {
		components := strings.Split(headContent, ": ")

		if !(len(components) >= expectedComponentes) {
			return "", fmt.Errorf("Oops! No hash available.")
		}

		branchPath := components[1]
		fullBranchPath := filepath.Join(pwd, git, branchPath)

		bytes, err = os.ReadFile(fullBranchPath)

		if err != nil {
			return "", fmt.Errorf("Oops! Something happened! %v", err)
		}

		return string(bytes), nil
	}

	/// HEAD points to a commit in particular.
	if len(headContent) > 0 {
		return headContent, nil
	}

	return "", fmt.Errorf("Oops! No hash available.")
}
