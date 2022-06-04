package cli

import "context"

// HelpCommand is in charge of printing the help text for a given command.
// its flags and any other information available to make it easy for the user.
type HelpCommand struct {
	commands []Command
}

func (c HelpCommand) Name() string {
	return "help"
}

func (c HelpCommand) HelpText() string {
	return "Provides help for a given command, p.e. depbot help list."
}

func (c HelpCommand) Main(ctx context.Context, pwd string, args []string) error {
	return nil
}
