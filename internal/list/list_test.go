package list_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/list"
)

func TestCommand(t *testing.T) {
	deps := []depbot.Dependency{
		{Name: "github.com/wawandco/ox", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
		{Name: "github.com/wawandco/maildoor", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
		{Name: "github.com/wawandco/fako", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
	}

	fakeFinder := func(wd string) (depbot.Dependencies, error) {
		return deps, nil
	}

	t.Run("No dependency found", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		c := &list.Command{}
		c.SetIO(out, out, nil)

		err := c.Main(context.Background(), t.TempDir(), []string{})
		if err != nil {
			t.Fatalf("error running find command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Total dependencies found: 0")) {
			t.Errorf("expected output to contain 'Total dependencies found:'")
		}
	})

	t.Run("One finder dep", func(t *testing.T) {
		c := list.NewCommand(fakeFinder)
		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)

		err := c.Main(context.Background(), t.TempDir(), []string{})
		if err != nil {
			t.Fatalf("error running list command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Total dependencies found: 3")) {
			t.Errorf("expected output to contain 'Total dependencies found: 3'")
		}
	})

	t.Run("Multiple finders", func(t *testing.T) {
		c := list.NewCommand(
			fakeFinder,
			fakeFinder,
		)

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)

		err := c.Main(context.Background(), t.TempDir(), []string{})
		if err != nil {
			t.Fatalf("error running list command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Total dependencies found: 6")) {
			t.Errorf("expected output to contain 'Total dependencies found: 6'")
		}
	})

	t.Run("table with the dependencies", func(t *testing.T) {
		c := list.NewCommand(
			fakeFinder,
		)

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)

		err := c.Main(context.Background(), t.TempDir(), []string{})
		if err != nil {
			t.Fatalf("error running list command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Total dependencies found: 3")) {
			t.Fatalf("expected output to contain 'Total dependencies found: 3'")
		}

		deps, err := fakeFinder("")
		for _, v := range deps {
			if !bytes.Contains(out.Bytes(), []byte(v.Name)) {
				t.Fatalf("expected output to contain %v", v.Name)
			}

			if !bytes.Contains(out.Bytes(), []byte(v.Version)) {
				t.Fatalf("expected output to contain %v", v.Version)
			}

			if !bytes.Contains(out.Bytes(), []byte(v.Version)) {
				t.Fatalf("expected output to contain %v", v.Direct)
			}
		}

	})

	t.Run("finder with output flag", func(t *testing.T) {
		c := list.NewCommand(
			fakeFinder,
		)

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)
		c.ParseFlags([]string{"--output=json"})

		err := c.Main(context.Background(), t.TempDir(), []string{})
		if err != nil {
			t.Fatalf("error running list command: %v", err)
		}

		jm, err := json.Marshal(deps)
		if err != nil {
			t.Fatalf("error marshal deps: %v", err)
		}

		if !strings.Contains(out.String(), string(jm)) {
			t.Fatalf("expected output to contain %v", string(jm))
		}

	})

	t.Run("finder with output csv", func(t *testing.T) {
		c := list.NewCommand(
			fakeFinder,
		)

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)
		c.ParseFlags([]string{"--output", "csv"})

		err := c.Main(context.Background(), t.TempDir(), []string{})
		if err != nil {
			t.Fatalf("error running list command: %v", err)
		}

		for _, v := range deps {
			line := fmt.Sprintf("\"%v\",\"%v\",\"%v\",\"%v\"\n", v.Name, v.Version, v.File, v.Direct)
			if !strings.Contains(out.String(), line) {
				t.Fatalf("expected output to contain %v", line)
			}
		}
	})
}
