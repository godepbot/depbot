package list_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/list"
)

func TestCommand(t *testing.T) {

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

		fakeFinder := func(wd string) ([]depbot.Dependency, error) {
			dd := []depbot.Dependency{
				{Name: "github.com/wawandco/ox", Version: "v1.0.0"},
				{Name: "github.com/wawandco/maildoor", Version: "v1.0.0"},
				{Name: "github.com/wawandco/fako", Version: "v1.0.0"},
			}

			return dd, nil
		}

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
		fakeFinder := func(wd string) ([]depbot.Dependency, error) {
			dd := []depbot.Dependency{
				{Name: "github.com/wawandco/ox", Version: "v1.0.0"},
				{Name: "github.com/wawandco/maildoor", Version: "v1.0.0"},
				{Name: "github.com/wawandco/fako", Version: "v1.0.0"},
			}

			return dd, nil
		}

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

}
