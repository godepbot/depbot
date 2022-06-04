package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/godepbot/depbot/cmd/cli"
)

type testCommand string

func (t testCommand) Name() string {
	return string(t)
}

func (t testCommand) HelpText() string {
	return fmt.Sprintf("runs the %v thing basically", t)
}

func (t testCommand) Main(ctx context.Context, pwd string, args []string) error {
	return nil
}

func TestUsage(t *testing.T) {
	out := bytes.NewBuffer([]byte{})

	commands := []cli.Command{
		testCommand("test"),
		testCommand("other"),
	}

	app := cli.NewApp(commands...)

	app.IO = cli.IO{
		Out: out,
		Err: out,
	}

	err := app.Usage()
	if err != nil {
		t.Fatalf("error running usage: %v", err)
	}

	if !bytes.Contains(out.Bytes(), []byte("Usage: depbot [command] [options]")) {
		t.Fatalf("expected output to contain 'Usage: depbot [command] [options]'")
	}

	if !bytes.Contains(out.Bytes(), []byte("Commands")) {
		t.Fatalf("expected output to contain 'Commands")
	}

	for _, v := range commands {
		if !bytes.Contains(out.Bytes(), []byte(v.Name())) {
			t.Fatalf("expected output to contain '%v'", v)
		}

		if ht, ok := v.(cli.HelpTexter); ok && !bytes.Contains(out.Bytes(), []byte(ht.HelpText())) {
			t.Fatalf("expected output to contain '%v'", ht.HelpText())
		}
	}

	// Making sure the help text is there.
	h := &cli.HelpCommand{}
	if !bytes.Contains(out.Bytes(), []byte(h.Name())) {
		t.Fatalf("expected output to contain '%v'", h.Name())
	}

	if !bytes.Contains(out.Bytes(), []byte(h.HelpText())) {
		t.Fatalf("expected output to contain '%v'", h.HelpText())
	}

}
