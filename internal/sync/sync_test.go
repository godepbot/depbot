package sync_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/sync"
)

type fakeServer struct {
	responseCode int

	receivedRequest *http.Request
}

func (fk *fakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(fk.responseCode)

	fk.receivedRequest = r
}

func TestSyncCommand(t *testing.T) {
	fakeFinder := func(wd string) (depbot.Dependencies, error) {
		dd := []depbot.Dependency{
			{Name: "github.com/wawandco/ox", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
			{Name: "github.com/wawandco/maildoor", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
			{Name: "github.com/wawandco/fako", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
		}

		return dd, nil
	}

	fakeRevisionFinder := func(string) (string, error) {
		return "0de1b3e18d9cd5cd96b12e608ca47eff046fed0a", nil
	}

	fkServer := fakeServer{
		responseCode: http.StatusOK,
	}

	server := httptest.NewServer(&fkServer)
	defer server.Close()

	t.Run("No dependency found to sync", func(t *testing.T) {
		c := &sync.Command{}

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)
		c.SetClient(server.Client())

		err := c.Main(context.Background(), "", []string{})
		if err == nil && (err != depbot.ErrorNoDependenciesFound) {
			t.Errorf("expected output to contain '%v'", depbot.ErrorNoDependenciesFound)
		}
	})

	t.Run("One sync dep", func(t *testing.T) {
		os.Setenv("DEPBOT_API_KEY", "SETWITHENV")

		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)
		c.WithRevisionFinder(fakeRevisionFinder)

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())
		c.ParseFlags([]string{"--server-address", server.URL})

		err := c.Main(context.Background(), "", []string{})

		if err != nil {
			t.Fatalf("error running sync command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("3 dependencies synchronized.")) {
			t.Errorf("expected output to contain '%v'", "3 dependencies synchronized.")
		}

		if fkServer.receivedRequest.Header.Get("Authorization") != "Bearer SETWITHENV" {
			t.Errorf("expected output to contain '%v'", "SETWITHENV")
		}

		hash, _ := fakeRevisionFinder("")
		if fkServer.receivedRequest.Header.Get("X-Revision-Hash") != hash {
			t.Errorf("expected output to contain '%v'", hash)
		}

	})

	t.Run("Multiple finders", func(t *testing.T) {
		os.Setenv("DEPBOT_API_KEY", "MULTIPLEAPIKEY")

		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(
			fakeFinder,
			fakeFinder,
		)
		c.WithRevisionFinder(fakeRevisionFinder)

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())
		c.ParseFlags([]string{"--server-address", server.URL})

		err := c.Main(context.Background(), "", []string{})

		if err != nil {
			t.Fatalf("error running sync command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("6 dependencies synchronized.")) {
			t.Errorf("expected output to contain '%v'", "dependencies synchronized.")
		}

		if fkServer.receivedRequest.Header.Get("Authorization") != "Bearer MULTIPLEAPIKEY" {
			t.Errorf("expected output to contain '%v'", "MULTIPLEAPIKEY")
		}

		hash, _ := fakeRevisionFinder("")
		if fkServer.receivedRequest.Header.Get("X-Revision-Hash") != hash {
			t.Errorf("expected output to contain '%v'", hash)
		}
	})

	t.Run("No API KEY", func(t *testing.T) {
		os.Setenv("DEPBOT_API_KEY", "")

		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())
		c.WithRevisionFinder(fakeRevisionFinder)

		c.ParseFlags([]string{"--server-address", server.URL})

		err := c.Main(context.Background(), "", []string{})
		if err == nil && (err != sync.ErrorMissingApiKey) {
			t.Errorf("expected output to contain '%v'", sync.ErrorMissingApiKey)
		}
	})

	t.Run("API Key passed as flag", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())
		c.WithRevisionFinder(fakeRevisionFinder)

		c.ParseFlags([]string{"--server-address", server.URL, "--api-key", "SETWITHFLAG"})

		err := c.Main(context.Background(), "", []string{})

		if err != nil {
			t.Fatalf("error running sync command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("3 dependencies synchronized.")) {
			t.Errorf("expected output to contain '%v'", "dependencies synchronized.")
		}

		if fkServer.receivedRequest.Header.Get("Authorization") != "Bearer SETWITHFLAG" {
			t.Errorf("expected output to contain '%v'", "SETWITHFLAG")
		}

		hash, _ := fakeRevisionFinder("")
		if fkServer.receivedRequest.Header.Get("X-Revision-Hash") != hash {
			t.Errorf("expected output to contain '%v'", hash)
		}
	})

	t.Run("Server Error", func(t *testing.T) {
		fkServer.responseCode = http.StatusInternalServerError

		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())
		c.WithRevisionFinder(fakeRevisionFinder)
		c.ParseFlags([]string{"--server-address", server.URL, "--api-key", "SETWITHFLAG"})

		err := c.Main(context.Background(), "", []string{})
		if err == nil {
			t.Fatalf("expected the sync command to error")
		}
	})

	t.Run("ENV and FLAG set", func(t *testing.T) {
		fkServer.responseCode = http.StatusOK
		os.Setenv("DEPBOT_API_KEY", "ENV")

		out := bytes.NewBuffer([]byte{})
		c := sync.NewCommand(fakeFinder)

		c.SetIO(out, out, nil)
		c.SetClient(server.Client())
		c.WithRevisionFinder(fakeRevisionFinder)
		c.ParseFlags([]string{"--server-address", server.URL, "--api-key", "FLAG"})

		err := c.Main(context.Background(), "", []string{})
		if err != nil {
			t.Fatalf("error running sync command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("3 dependencies synchronized.")) {
			t.Errorf("expected output to contain '%v'", "dependencies synchronized.")
		}

		if fkServer.receivedRequest.Header.Get("Authorization") != "Bearer FLAG" {
			t.Errorf("expected output to contain '%v'", "FLAG")
		}
	})

}
