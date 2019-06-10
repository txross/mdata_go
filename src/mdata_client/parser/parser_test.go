package parser

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/create"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/delete"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/list"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/set"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/show"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/update"
	"testing"
)

func TestGetParser(t *testing.T) {
	tests := map[string]struct {
		inCommand        []commands.Command
		inActiveCommand  []string
		outParser        *flags.Parser
		outActiveCommand string
	}{
		"noCommands": {
			inCommand:        []commands.Command{},
			inActiveCommand:  []string{"create", "12345678901234"},
			outParser:        flags.NewNamedParser("mdata", flags.Default),
			outActiveCommand: "",
		},
		"oneCommand": {
			inCommand:        []commands.Command{&create.Create{}},
			inActiveCommand:  []string{"create", "12345678901234"},
			outParser:        flags.NewNamedParser("mdata", flags.Default),
			outActiveCommand: "create",
		},
		"multipleCommands": {
			inCommand: []commands.Command{&create.Create{},
				&delete.Delete{},
				&update.Update{},
				&set.Set{},
				&show.Show{},
				&list.List{}},
			inActiveCommand:  []string{"list"},
			outParser:        flags.NewNamedParser("mdata", flags.Default),
			outActiveCommand: "list",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		parser := GetParser(test.inCommand)
		args, _ := parser.ParseArgs(test.inActiveCommand)

		// Test commands registered to parser
		assert.Equal(t, test.outParser, parser)

		// Test active command
		assert.Equal(t, test.outActiveCommand, parser.Active.Name)
	}
}
