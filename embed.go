package depbot

import (
	"embed"

	"github.com/paganotoni/fsbox"
)

var (
	//go:embed app/templates public migrations config
	fs embed.FS

	// Boxes used by the app, these are based on the embed.FS declared
	// in the fs variable.
	Assets     = fsbox.New(fs, "public")
	Templates  = fsbox.New(fs, "app/templates")
	Migrations = fsbox.New(fs, "migrations")
	Config     = fsbox.New(fs, "config")
)
