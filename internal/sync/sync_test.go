package sync_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/sync"
)

func mockEndPoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(sync.ErrorNoSyncDep.Error()))
	}

	if strings.Contains(string(body), "contains something wrong") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("bad request because i want to be like it"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("dependencies synchronized."))
	}
}

func TestSyncCommand(t *testing.T) {

	dir := t.TempDir()

	fakeFinder := func(wd string) (depbot.Dependencies, error) {
		dd := []depbot.Dependency{
			{Name: "github.com/wawandco/ox", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
			{Name: "github.com/wawandco/maildoor", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
			{Name: "github.com/wawandco/fako", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
		}

		return dd, nil
	}

	gitDir := filepath.Join(dir, ".git")
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

	server := httptest.NewServer(http.HandlerFunc(mockEndPoint))
	defer server.Close()

	os.Setenv(sync.DepbotServerAddr, server.URL)
	os.Setenv(sync.DepbotApiKey, "An API Key")

	t.Run("No dependency found to sync", func(t *testing.T) {

		out := bytes.NewBuffer([]byte{})
		c := &sync.Command{}

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})
		if err == nil && (err != depbot.ErrorNoDependenciesFound) {
			t.Errorf("expected output to contain '%v'", depbot.ErrorNoDependenciesFound)
		}
	})

	t.Run("One sync dep", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)
		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})

		if err != nil {
			t.Fatalf("error running sync command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("dependencies synchronized.")) {
			t.Errorf("expected output to contain '%v'", "dependencies synchronized.")
		}

		if !bytes.Contains(out.Bytes(), []byte("3")) {
			t.Errorf("expected output to contain '%v'", 3)
		}
	})

	t.Run("Multiple finders", func(t *testing.T) {
		c := sync.NewCommand(
			fakeFinder,
			fakeFinder,
		)

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})
		if err != nil {
			t.Fatalf("error running sync command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("dependencies synchronized.")) {
			t.Errorf("expected output to contain '%v'", "dependencies synchronized.")
		}

		if !bytes.Contains(out.Bytes(), []byte("6")) {
			t.Errorf("expected output to contain '%v'", 6)
		}
	})

	t.Run("Bad Resnponse", func(t *testing.T) {
		fakeBadFinder := func(wd string) (depbot.Dependencies, error) {
			dd := []depbot.Dependency{
				{Name: "contains something wrong"},
			}

			return dd, nil
		}
		c := sync.NewCommand(
			fakeBadFinder,
		)

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})
		if err == nil && (err != sync.ErrorNoSyncDep) {
			t.Errorf("expected output to contain '%v'", sync.ErrorNoSyncDep)
		}
	})

	t.Run("Error Api key", func(t *testing.T) {
		os.Setenv(sync.DepbotApiKey, "")
		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)
		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})
		if err == nil && (err != sync.ErrorMissingApiKey) {
			t.Errorf("expected output to contain '%v'", sync.ErrorMissingApiKey)
		}
	})

	t.Run("Sync command with args", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)
		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})
		if err == nil && (err != sync.ErrorMissingApiKey) {
			t.Errorf("expected output to contain '%v'", sync.ErrorMissingApiKey)
		}

		err = c.Main(context.Background(), dir, []string{"--api-key=API_KEY"})
		if err != nil {
			t.Errorf("expected output to no contain errors")
		}

		if os.Getenv(sync.DepbotApiKey) != "API_KEY" {
			t.Errorf("expected env variable to be 'API_KEY'")
		}

		c.Main(context.Background(), dir, []string{"--api-key=Other_Key", "--server-address=my.server.com"})
		if os.Getenv(sync.DepbotApiKey) != "Other_Key" {
			t.Errorf("expected env variable to be 'Other_Key' got %v instead", os.Getenv(sync.DepbotApiKey))
		}
		if os.Getenv(sync.DepbotServerAddr) != "my.server.com" {
			t.Errorf("expected env variable to be 'my.server.com' got %v instead", os.Getenv(sync.DepbotServerAddr))
		}

		c.Main(context.Background(), dir, []string{"--api-key=Other_Key", fmt.Sprintf("--server-address=%v", server.URL)})
		if os.Getenv(sync.DepbotApiKey) != "Other_Key" {
			t.Errorf("expected env variable to be 'Other_Key' got %v instead", os.Getenv(sync.DepbotApiKey))
		}
		if os.Getenv(sync.DepbotServerAddr) != server.URL {
			t.Errorf("expected env variable to be '%v' got %v", server.URL, os.Getenv(sync.DepbotServerAddr))
		}

	})

}
