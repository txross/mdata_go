package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/create"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/list"
	"testing"
)

func TestGetParser(t *testing.T) {
	tests := map[string]struct {
		inCommand        []commands.Command
		inActiveCommand  []string
		outActiveCommand string
	}{
		"oneCommand": {
			inCommand:        []commands.Command{&create.Create{}},
			inActiveCommand:  []string{"create", "12345678901234"},
			outActiveCommand: "create",
		},
		"multipleCommands": {
			inCommand: []commands.Command{&create.Create{},
				&list.List{}},
			inActiveCommand:  []string{"list"},
			outActiveCommand: "list",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		parser := GetParser(test.inCommand)
		parser.ParseArgs(test.inActiveCommand)

		// Test active command
		assert.Equal(t, test.outActiveCommand, parser.Active.Name)
	}
}
