package sync_test

import (
	"bytes"
	"context"
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
		w.Write([]byte(depbot.MessageError_NoSyncDep))
	}

	if strings.Contains(string(body), "contains something wrong") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("bad request because i want to be like it"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(depbot.MessageSucces_SyncDep))
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

	os.Setenv(depbot.EnvVariable_ServerADDR, server.URL)
	os.Setenv(depbot.EnvVariable_ApiKey, "An API Key")

	t.Run("No dependency found to sync", func(t *testing.T) {

		out := bytes.NewBuffer([]byte{})
		c := &sync.Command{}

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})
		if err == nil {
			t.Fatalf("expected error to contain: '%v'", depbot.MessageError_NoDependencies)
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

		if !bytes.Contains(out.Bytes(), []byte(depbot.MessageSucces_SyncDep)) {
			t.Errorf("expected output to contain '%v'", depbot.MessageSucces_SyncDep)
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

		if !bytes.Contains(out.Bytes(), []byte(depbot.MessageSucces_SyncDep)) {
			t.Errorf("expected output to contain '%v'", depbot.MessageSucces_SyncDep)
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
		if err == nil && !strings.Contains(err.Error(), depbot.MessageError_NoSyncDep) {
			t.Errorf("expected output to contain '%v'", depbot.MessageError_NoSyncDep)
		}
	})

	t.Run("Error Api key", func(t *testing.T) {
		os.Setenv(depbot.EnvVariable_ApiKey, "")
		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)
		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), dir, []string{})

		if err == nil {
			t.Errorf("expected output to contain '%v'", depbot.MessageError_MissingApiKey)
		}
	})

}
