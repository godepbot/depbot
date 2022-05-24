package list_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/list"
)

func TestCommand(t *testing.T) {

	fakeFinder := func(wd string) (depbot.DependencyAnalisys, error) {
		dd := []depbot.Dependency{
			{Name: "github.com/wawandco/ox", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
			{Name: "github.com/wawandco/maildoor", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
			{Name: "github.com/wawandco/fako", Version: "v1.0.0", Kind: depbot.DependencyKindLibrary, File: "go.mod", Direct: true},
		}

		da := depbot.DependencyAnalisys{
			Timestamp:    time.Now().Unix(),
			Dependencies: dd,
		}

		return da, nil
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

		da, err := fakeFinder("")
		for _, v := range da.Dependencies {
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

}
