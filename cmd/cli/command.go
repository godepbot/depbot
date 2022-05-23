package cli

import (
	"context"
)

// Command is a struct that can be invoked by the main app
// to do so it would use the Name method to identify it.
type Command interface {
	Name() string
	Main(ctx context.Context, pwd string, args []string) error
}
