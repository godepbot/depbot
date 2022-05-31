package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	headContent = "ref: refs/heads/test-branch-1234"
	headsPath   = "refs/heads"
	branchName  = "test-branch-1234"
	hash        = "0de1b3e18d9cd5cd96b12e608ca47eff046fed0a"
)

func Test_HeadLinkedToBranch(t *testing.T) {

	fmt.Println("Executing this test!")
	tmp := t.TempDir()
	gitDir := filepath.Join(tmp, ".git")
	err := os.Mkdir(gitDir, os.ModePerm)

	if err != nil {
		fmt.Println("Error creating directory is:", err)
		return
	}

	headPath := filepath.Join(gitDir, "HEAD")

	err = os.WriteFile(headPath, []byte(headContent), os.ModePerm)

	if err != nil {
		fmt.Println("Error writing file is:", err)
		return
	}

	headsDir := filepath.Join(gitDir, headsPath)

	err = os.MkdirAll(headsDir, os.ModePerm)

	if err != nil {
		fmt.Println("Error creating directory is:", err)
		return
	}

	branchFile := filepath.Join(headsDir, branchName)

	os.WriteFile(branchFile, []byte(hash), os.ModePerm)

	data, err := FindLatestHash(tmp)

	if data != hash {
		t.Fatalf("Got %v, but was expected %v", data, hash)
	}

}

func Test_NoGitRepository(t *testing.T) {
	tmp := t.TempDir()

	_, err := FindLatestHash(tmp)

	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("Result was: %v, but was expected an error including -> %v", err.Error(), "no such file or directory")
	}

}
