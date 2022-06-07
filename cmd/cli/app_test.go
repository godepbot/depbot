package cli_test

import (
	"context"
	"fmt"
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
