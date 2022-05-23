package cli

import (
	"context"
	"fmt"
	"io"
)

type App struct {
	IO

	Commands []Command
}

func (app *App) findCommand(args []string) Command {
	for _, v := range app.Commands {
		if v.Name() == args[0] {
			return v
		}
	}

	return nil
}

func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	if len(args) > 0 {
		return app.Usage(app.Stdout())
	}

	command := app.findCommand(args)
	if command == nil {
		return app.Usage(app.Stdout())
	}

	if ist, ok := command.(IOSetter); ok {
		ist.SetIO(app.Stdout(), app.Stderr(), app.Stdin())
	}

	return command.Main(ctx, pwd, args[1:])
}

func (app *App) Usage(w io.Writer) error {
	fmt.Fprintln(w, "Usage: depbot [command] [options]")
	fmt.Fprintln(w, "---------------")

	return nil
}
