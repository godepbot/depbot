package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	gitPath             = ".git/"
	head                = "HEAD"
	expectedComponentes = 2
	git                 = ".git"
	refsPath            = "refs/"
)

func FindLatestHash(pwd string) (string, error) {
	path := filepath.Join(pwd, git, head)
	bytes, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Error is :)", err)
		return "", err
	}

	headContent := strings.TrimSpace(string(bytes))

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

	if len(headContent) > 0 {
		return headContent, nil
	}

	return "", fmt.Errorf("Oops! No hash available.")

}
