package cli

import (
	"context"
	"fmt"
	"io"
)

type App struct {
	IO
}

func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	if len(args) > 0 {
		return app.Usage(app.Stdout())
	}

	// if app.Commands == nil {
	// 	app.Commands = map[string]Commander{}
	// }

	// cmd, ok := app.Commands[args[0]]
	// if !ok {
	// 	return fmt.Errorf("command %q not found", args[0])
	// }

	return nil //cmd.Main(ctx, pwd, args[1:])
}

func (app *App) Usage(w io.Writer) error {
	fmt.Fprintln(w, "Usage: depbot [options]")
	fmt.Fprintln(w, "---------------")

	return nil
}
