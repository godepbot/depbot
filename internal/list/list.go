package list

import (
	"context"
	"fmt"
	"io"

	"github.com/godepbot/depbot"
)

type Command struct {
	finders []depbot.FinderFn

	stderr io.Writer
	stdout io.Writer
	stdin  io.Reader
}

func (c *Command) Name() string {
	return "find"
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
