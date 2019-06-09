package commands

import flags "github.com/jessevdk/go-flags"

// All subcommands implement this interface
type Command interface {
	Register(*flags.Command) error
	Name() string
	KeyfilePassed() string
	UrlPassed() string
	Run() ([]byte, error)
}
