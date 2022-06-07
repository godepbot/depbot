package cli

import (
	"bytes"
	"context"
	"flag"
	"fmt"
)

// HelpCommand is in charge of printing the help text for a given command.
// its flags and any other information available to make it easy for the user.
type HelpCommand struct {
	IO
	flagSet *flag.FlagSet
	output  string

	Commands Commands
}

func (c HelpCommand) Name() string {
	return "help"
}

func (c HelpCommand) HelpText() string {
	return "Provides help for a given command, p.e. depbot help list."
}

func (c HelpCommand) ParseFlags(args []string) (*flag.FlagSet, error) {
	c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.flagSet.StringVar(&c.output, "output", "plain", "Output format. Can be plain, json or csv.")

	// This is to keep it silent
	c.flagSet.SetOutput(bytes.NewBuffer([]byte{}))
	c.flagSet.Usage = func() {}

	// Ignore the error we don't care if any error happens while parsing.
	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c HelpCommand) Main(ctx context.Context, pwd string, args []string) error {
	var command Command
	if len(args) > 0 {
		command = c.Commands.Find(args[0])
		if command == nil {
			fmt.Fprintf(c.Stdout(), "Error: did not find `%v` command\n", args[0])
		}
	}

	if len(args) == 0 || command == nil {
		return c.general()
	}

	return c.specific(command)
}

//general method prints the
func (c HelpCommand) general() error {
	fmt.Fprint(c.Stdout(), "Usage: depbot [command] [options]\n\n")

	// If there are no commands it just prints the usage.
	if len(c.Commands) == 0 {
		return nil
	}

	fmt.Fprintln(c.Stdout(), "Available Commands")
	fmt.Fprintln(c.Stdout(), "------------------")
	for _, v := range c.Commands {
		if ht, ok := v.(HelpTexter); ok {
			fmt.Fprintf(c.Stdout(), "%v\t%v\n", v.Name(), ht.HelpText())
			continue
		}

		fmt.Fprintf(c.Stdout(), "%v\t (runs the %[1]v command)\n", v.Name())
	}

	fmt.Fprintln(c.Stdout(), "\nFor command specific information use the help command, p.e.")
	fmt.Fprintln(c.Stdout(), "$ depbot help [command]")

	return nil
}

func (c HelpCommand) specific(cm Command) error {
	fmt.Fprintf(c.Stdout(), "Usage: depbot %v [options]\n\n", cm.Name())

	if ht, ok := cm.(HelpTexter); ok {
		fmt.Fprintf(c.Stdout(), ht.HelpText()+"\n\n")
	}

	if fl, ok := cm.(FlagParser); ok {
		fl, _ := fl.ParseFlags([]string{})
		fmt.Fprintf(c.Stdout(), "Flags:\n")
		fl.VisitAll(func(ff *flag.Flag) {
			fmt.Fprintf(c.Stdout(), "--%v\t%v\n", ff.Name, ff.Usage)
		})
	}

	return nil
}
