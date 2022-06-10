package list

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/godepbot/depbot"
)

type Command struct {
	finders []depbot.FinderFn
	flagSet *flag.FlagSet

	stderr io.Writer
	stdout io.Writer
	stdin  io.Reader

	output string
}

func (c *Command) Name() string {
	return "list"
}

func (c *Command) HelpText() string {
	return "Analyzes and lists dependencies by walking the current directory for dependency files."
}

func (c *Command) Main(ctx context.Context, pwd string, args []string) error {
	deps := []depbot.Dependency{}
	for _, df := range c.finders {
		dx, err := df(pwd)
		if err != nil {
			return err
		}

		deps = append(deps, dx...)
	}

	w := new(tabwriter.Writer)
	w.Init(c.stdout, 0, 8, 0, '\t', 0)

	switch c.output {
	case "json":
		jm, err := json.Marshal(deps)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%v", string(jm))
	case "csv":
		fmt.Fprintln(w, "\"Name\",\"Version\",\"File\",\"Direct\"")
		for _, v := range deps {
			fmt.Fprintf(w, "\"%v\",\"%v\",\"%v\",\"%v\"\n", v.Name, v.Version, v.File, v.Direct)
		}
	default:
		fmt.Fprintln(c.stdout, "Total dependencies found:", len(deps))
		fmt.Fprintln(w, "\nName\tVersion\tFile\tDirect")
		fmt.Fprintln(w, "----\t-------\t----\t-------")
		for _, v := range deps {
			fmt.Fprintf(
				w,
				"%v\t%v\t%v\t%v\t\n",
				v.Name,
				v.Version,
				v.File,
				v.Direct,
			)
		}
	}

	fmt.Fprintln(w)
	w.Flush()

	return nil
}

func (c *Command) SetIO(stderr io.Writer, stdout io.Writer, stdin io.Reader) {
	c.stderr = stderr
	c.stdout = stdout
	c.stdin = stdin
}

// NewCommand with the given finder function.
func NewCommand(finders ...depbot.FinderFn) *Command {
	c := &Command{
		finders: finders,
	}

	fls := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	fls.StringVar(&c.output, "output", "plain", "Output format. Can be plain, json or csv.")

	// This is to keep it silent
	fls.SetOutput(bytes.NewBuffer([]byte{}))
	fls.Usage = func() {}

	c.flagSet = fls

	return c
}

func (c *Command) ParseFlags(args []string) (*flag.FlagSet, error) {
	// Ignore the error we don't care if any error happens while parsing.
	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}
