package list

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/godepbot/depbot"
)

type Command struct {
	finders []depbot.FinderFn

	stderr io.Writer
	stdout io.Writer
	stdin  io.Reader
}

func (c *Command) Name() string {
	return "list"
}

func (c *Command) Main(ctx context.Context, pwd string, args []string) error {
	deps := []depbot.Dependency{}
	for _, df := range c.finders {
		dx, err := df(pwd)
		if err != nil {
			return err
		}

		deps = append(deps, dx...)
	}

	fmt.Fprintln(c.stdout, "Total dependencies found:", len(deps))

	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(c.stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "\nName\tVersion\tFile\tDirect")
	fmt.Fprintln(w, "----\t-------\t----\t-------")
	for _, v := range deps {
		fmt.Fprintf(
			w,
			"%v\t%v\t%v\t%v\t\n",
			v.Name,
			v.Version,
			v.File,
			v.Direct,
		)
	}
	fmt.Fprintln(w)
	w.Flush()

	return nil
}

func (c *Command) SetIO(stderr io.Writer, stdout io.Writer, stdin io.Reader) {
	c.stderr = stderr
	c.stdout = stdout
	c.stdin = stdin
}

// NewCommand with the given finder function.
func NewCommand(finders ...depbot.FinderFn) *Command {
	return &Command{
		finders: finders,
	}
}
