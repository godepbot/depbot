package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	commitHash = "0de1b3e18d9cd5cd96b12e608ca47eff046fed0a"
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
	headContent := "ref: refs/heads/test-branch-1234"
	err = os.WriteFile(headPath, []byte(headContent), os.ModePerm)
	if err != nil {
		fmt.Println("Error writing file is:", err)
		return
	}

	headsPath := "refs/heads"
	headsDir := filepath.Join(gitDir, headsPath)

	err = os.MkdirAll(headsDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory is:", err)
		return
	}

	branchName := "test-branch-1234"
	branchFile := filepath.Join(headsDir, branchName)
	os.WriteFile(branchFile, []byte(commitHash), os.ModePerm)

	data, err := FindLatestHash(tmp)

	if data != commitHash {
		t.Fatalf("Got %v, but was expected %v", data, commitHash)
	}
}

func Test_NoGitRepository(t *testing.T) {
	tmp := t.TempDir()

	_, err := FindLatestHash(tmp)
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("Result was: %v, but was expected an error including -> %v", err.Error(), "no such file or directory")
	}

}

func Test_HeadLinkedToCommit(t *testing.T) {
	tmp := t.TempDir()

	gitDir := filepath.Join(tmp, ".git")
	err := os.Mkdir(gitDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory is:", err)
		return
	}

	headPath := filepath.Join(gitDir, "HEAD")
	err = os.WriteFile(headPath, []byte(commitHash), os.ModePerm)
	if err != nil {
		fmt.Println("Error is h", err)
	}

	data, _ := FindLatestHash(tmp)
	if data != commitHash {
		t.Fatalf("Result was: %v, but was expected -> %v", data, commitHash)
	}

}
