package cli

import (
	"context"
	"fmt"
)

type App struct {
	IO

	commands []Command
}

// findCommand with given args, if no args passed
// it will return nil.
func (app *App) findCommand(args []string) Command {
	if len(args) == 0 {
		return nil
	}

	for _, v := range app.commands {
		if v.Name() == args[0] {
			return v
		}
	}

	return nil
}

// Main entry point for the application. This method finds the passed command
// and executes it with the passed arguments. If there is no command passed
// it will print the usage.
func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	if len(args) == 0 {
		return app.Usage()
	}

	command := app.findCommand(args)
	if command == nil {
		return app.Usage()
	}

	if ist, ok := command.(IOSetter); ok {
		ist.SetIO(app.Stdout(), app.Stderr(), app.Stdin())
	}

	return command.Main(ctx, pwd, args[1:])
}

// Usage of the App, it will print basic usage information
// and a list of commands available.
func (app *App) Usage() error {
	fmt.Fprint(app.Stdout(), "Usage: depbot [command] [options]\n\n")

	// If there are no commands it just prints the usage.
	if len(app.commands) == 0 {
		return nil
	}

	fmt.Fprintln(app.Stdout(), "Commands")
	fmt.Fprintln(app.Stdout(), "------------------")
	for _, v := range app.commands {
		if ht, ok := v.(HelpTexter); ok {
			fmt.Fprintf(app.Stdout(), "%v\t%v\n", v.Name(), ht.HelpText())
			continue
		}

		fmt.Fprintf(app.Stdout(), "%v\t (runs the %[1]v command)\n", v.Name())
	}

	fmt.Fprintln(app.Stdout(), "\nFor command specific information use the help command, p.e.")
	fmt.Fprintln(app.Stdout(), "$ depbot help [command]")

	return nil
}

// NewApp creates a CLI app with the given commands.
// It prepends the `help` command to the list of commands.
func NewApp(commands ...Command) *App {
	// Adding the help command always.
	cmms := []Command{
		&HelpCommand{
			commands: commands,
		},
	}

	return &App{
		commands: append(cmms, commands...),
	}
}
