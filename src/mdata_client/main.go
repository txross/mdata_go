/**
 * Copyright 2018 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

package main

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	flags "github.com/jessevdk/go-flags"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/create"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/delete"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/list"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/set"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/show"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands/update"
	"github.com/tross-tyson/mdata_go/src/mdata_client/constants"
	"github.com/tross-tyson/mdata_go/src/mdata_client/rest_service"
	"os"
)

var DISTRIBUTION_VERSION string

var logger *logging.Logger = logging.Get()

func init() {
	if len(constants.DISTRIBUTION_VERSION) == 0 {
		DISTRIBUTION_VERSION = "Unknown"
	}
}

func runCommandLine(parser *flags.Parser, remaining []string) {

	Commands := []commands.Command{
		&create.Create{},
		&delete.Delete{},
		&update.Update{},
		&set.Set{},
		&show.Show{},
		&list.List{},
	}

	for _, cmd := range Commands {
		err := cmd.Register(parser.Command)
		if err != nil {
			logger.Errorf("Couldn't register command %v: %v", cmd.Name(), err)
			os.Exit(1)
		}
	}

	_, err := parser.ParseArgs(remaining)

	fmt.Printf("ALL COMMAND LINE ARGUMENTS: \n\t%v", parser.Command.Active)

	if err != nil {
		logger.Errorf("Error parsing commands %v: %v", remaining, err)
		os.Exit(1)
	}

	// If a sub-command was passed, run it
	if parser.Command.Active == nil {
		os.Exit(2)
	}

	name := parser.Command.Active.Name
	for _, cmd := range Commands {
		if cmd.Name() == name {
			response, err := cmd.Run()
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			fmt.Println(response)
			return
		}
	}

	fmt.Println("Error: Command not found: ", name)

	return
}

type Opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Enable more verbose output"`
	Version bool   `short:"V" long:"version" description:"Display version information"`
	Server  bool   `short:"S" long:"server" description:"Run as REST Server instead of command line"`
	Port    uint   `short:"p" long:"port" description:"Provide the port to run the REST Service. Default -p=8888"`
}

func main() {
	arguments := os.Args[1:]
	for _, arg := range arguments {
		if arg == "-V" || arg == "--version" {
			fmt.Println(constants.DISTRIBUTION_NAME + " (Hyperledger Sawtooth) version " + constants.DISTRIBUTION_VERSION)
			os.Exit(0)
		}
	}

	var opts Opts
	parser := flags.NewParser(&opts, flags.Default)
	parser.Command.Name = "mdata"

	// Set verbosity
	switch len(opts.Verbose) {
	case 2:
		logger.SetLevel(logging.DEBUG)
	case 1:
		logger.SetLevel(logging.INFO)
	default:
		logger.SetLevel(logging.WARN)
	}

	remaining, err := parser.Parse()
	if e, ok := err.(*flags.Error); ok {
		if e.Type == flags.ErrHelp {
			return
		} else {
			os.Exit(1)
		}
	}

	if opts.Server {
		// Instantiate RESTful API
		rest_service.Run(opts.Port)
	} else {
		runCommandLine(parser, remaining)
	}

}
