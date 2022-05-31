package commit

import (
	"fmt"
	"os"
	"strings"
)

const (
	gitPath             = ".git/"
	headPath            = ".git/HEAD"
	expectedComponentes = 2
)

func FindLatestHash() (string, error) {

	bytes, err := os.ReadFile(headPath)

	if err != nil {
		return "", err
	}

	headContent := strings.TrimSpace(string(bytes))

	components := strings.Split(headContent, ": ")

	if !(len(components) >= expectedComponentes) {
		return "", fmt.Errorf("Oops! Something happened!")
	}

	branchPath := components[1]
	fullBranchPath := fmt.Sprintf("%s%s", gitPath, branchPath)
	bytes, err = os.ReadFile(fullBranchPath)

	if err != nil {
		return "", fmt.Errorf("Oops! Something happened!")
	}

	return string(bytes), nil
}
