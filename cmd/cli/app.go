package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/gomodules"
	"github.com/godepbot/depbot/internal/jsmodules"
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

	deps, err := findDependencies(pwd)
	if err != nil {
		fmt.Println("Error is", err)
		return fmt.Errorf("error finding dependencies: %w", err)
	}

	fmt.Println("Total dependencies found:", len(deps))

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

func findDependencies(pwd string) (depbot.Dependencies, error) {
	deps := depbot.Dependencies{}

	goDeps, err := gomodules.FindDependencies(pwd)
	if err != nil {
		return deps, fmt.Errorf("error finding go dependencies, err : %w", err)
	}

	deps = append(deps, goDeps...)

	jsDeps, err := jsmodules.FindDependencies(pwd)
	if err != nil {
		return deps, fmt.Errorf("error finding js dependencies, err : %w", err)
	}

	deps = append(deps, jsDeps...)

	return deps, nil
}
