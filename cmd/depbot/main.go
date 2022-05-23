package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/godepbot/depbot/cmd/cli"
	"github.com/godepbot/depbot/internal/find"
)

// app for the CLI, commands used will be added here.
var app = &cli.App{
	Commands: []cli.Command{
		&find.Command{},
	},
}

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
