package find

import (
	"context"
	"fmt"
	"io"

	"github.com/godepbot/depbot/internal/gomodules"
)

type Command struct {
	stderr io.Writer
	stdout io.Writer
	stdin  io.Reader
}

func (c *Command) Name() string {
	return "find"
}

func (c *Command) Main(ctx context.Context, pwd string, args []string) error {
	deps, err := gomodules.FindDependencies(pwd)
	if err != nil {
		return fmt.Errorf("error finding dependencies: %w", err)
	}

	fmt.Fprintln(c.stdout, "Total dependencies found:", len(deps))

	return nil
}

func (c *Command) SetIO(stderr io.Writer, stdout io.Writer, stdin io.Reader) {
	c.stderr = stderr
	c.stdout = stdout
	c.stdin = stdin
}
