package parser

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/create"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/delete"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/list"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/set"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/show"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/update"
)

func Commands() []commands.Command {
	return []commands.Command{
		&create.Create{},
		&delete.Delete{},
		&update.Update{},
		&set.Set{},
		&show.Show{},
		&list.List{},
	}
}

func GetParser() *flags.Parser {
	var opts Opts
	p := flags.NewNamedParser("mdata", *flags.Default)

	for _, cmd := range Commands() {
		err := cmd.Register(p.Command)
		if err != nil {
			logger.Errorf("Couldn't register command %v: %v", cmd.Name(), err)
			os.Exit(1)
		}
	}

	return &p
}
