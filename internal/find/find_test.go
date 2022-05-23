package find_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/godepbot/depbot/internal/find"
)

func TestCommand(t *testing.T) {

	t.Run("No dependency files", func(t *testing.T) {
		c := &find.Command{}

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)
		wd := t.TempDir()

		err := c.Main(context.Background(), wd, []string{})
		if err != nil {
			t.Fatalf("error running find command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Total dependencies found: 0")) {
			t.Errorf("expected output to contain 'Total dependencies found:'")
		}
	})

	t.Run("Dependency File there", func(t *testing.T) {

		wd := t.TempDir()
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("error changing directory: %v", err)
		}

		err := os.WriteFile(
			filepath.Join(wd, "go.mod"),
			[]byte(`module something
			go 1.18

require (
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
)
			`),
			0777,
		)

		if err != nil {
			t.Fatalf("error writing go.mod file: %v", err)
		}

		c := &find.Command{}

		out := bytes.NewBuffer([]byte{})
		c.SetIO(out, out, nil)

		err = c.Main(context.Background(), wd, []string{})
		if err != nil {
			t.Fatalf("error running find command: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Total dependencies found: 3")) {
			t.Errorf("expected output to contain 'Total dependencies found:'")
		}
	})

}
