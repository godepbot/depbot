package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/godepbot/depbot/cmd/cli"
	"github.com/godepbot/depbot/internal/gomodules"
	"github.com/godepbot/depbot/internal/jspackages"
	"github.com/godepbot/depbot/internal/list"
)

// app for the CLI, commands used will be added here.
var (
	app = &cli.App{
		Commands: []cli.Command{
			// find command
			list.NewCommand(
				gomodules.FindDependencies,
				jspackages.FindPackageDependencies,
				jspackages.FindPackageLockDependencies,
			),
		},
	}
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = app.Main(ctx, pwd, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-ctx.Done()
}
