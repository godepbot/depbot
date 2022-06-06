package revision

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_HeadLinkedToBranch(t *testing.T) {
	tmp := t.TempDir()
	gitDir := filepath.Join(tmp, ".git")
	err := os.Mkdir(gitDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Error %v creating directory at %v", gitDir, err)
		return
	}

	headPath := filepath.Join(gitDir, "HEAD")
	err = os.WriteFile(headPath, []byte("ref: refs/heads/test-branch-1234"), os.ModePerm)
	if err != nil {
		t.Fatalf("Error writing file is: %v", err)
		return
	}

	headsDir := filepath.Join(gitDir, "refs/heads")

	err = os.MkdirAll(headsDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Error creating directory is: %v", err)
		return
	}

	commitHash := "0de1b3e18d9cd5cd96b12e608ca47eff046fed0a"
	branchFile := filepath.Join(headsDir, "test-branch-1234")
	os.WriteFile(branchFile, []byte(commitHash), os.ModePerm)

	data, _ := FindLatestHash(tmp)
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
		t.Fatalf("Error %v creating directory at %v", gitDir, err)
		return
	}

	commitHash := "3cb32c9ec6ae2f88dbee1e8be923a691dd73427c"
	headPath := filepath.Join(gitDir, "HEAD")
	err = os.WriteFile(headPath, []byte(commitHash), os.ModePerm)
	if err != nil {
		t.Fatalf("Error %v writing file at %v", err, headPath)
	}

	data, _ := FindLatestHash(tmp)
	if data != commitHash {
		t.Fatalf("Result was: %v, but was expected -> %v", data, commitHash)
	}
}

func Test_HeadLinkedToCommit_WithNewLineCharacter(t *testing.T) {
	tmp := t.TempDir()

	gitDir := filepath.Join(tmp, ".git")
	err := os.Mkdir(gitDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Error %v creating directory at %v", gitDir, err)
		return
	}

	commitHash := "3cb32c9ec6ae2f88dbee1e8be923a691dd73427c\n"
	headPath := filepath.Join(gitDir, "HEAD")
	err = os.WriteFile(headPath, []byte(commitHash), os.ModePerm)
	if err != nil {
		t.Fatalf("Error %v writing file at %v", err, headPath)
	}

	data, _ := FindLatestHash(tmp)
	if strings.Contains(data, "\n") {
		t.Errorf("result contain the new line character")
	}

	if data != strings.ReplaceAll(commitHash, "\n", "") {
		t.Fatalf("Result was: %v, but was expected -> %v", data, strings.ReplaceAll(commitHash, "\n", ""))
	}
}
