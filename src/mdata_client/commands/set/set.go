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

package set

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/tross-tyson/mdata_go/src/mdata_client/client"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands"
)

type Set struct {
	Args struct {
		Gtin  string `positional-arg-name:"gtin" required:"true" description:"Identify the gtin of the product to set state"`
		State string `positional-arg-name:"state" required:"true" description:"Specify the state to set the <gtin>: ACTIVE, INACTIVE, DISCONTINUED" choice:"INACTIVE" choice:"ACTIVE" choice:"DISCONTINUED"`
	} `positional-args:"true"`
	Url     string `long:"url" description:"Specify URL of REST API"`
	Keyfile string `long:"keyfile" description:"Identify file containing user's private key"`
	Wait    uint   `long:"wait" description:"Set time, in seconds, to wait for transaction to commit"`
}

func (args *Set) Name() string {
	return "set"
}

func (args *Set) KeyfilePassed() string {
	return args.Keyfile
}

func (args *Set) UrlPassed() string {
	return args.Url
}

func (args *Set) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Set state of a product", "Sends an mdata transaction to set state of <gtin> to <state>.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *Set) Run() (string, error) {
	// Construct client
	gtin := args.Args.Gtin
	state := args.Args.State
	wait := args.Wait

	mdataClient, err := client.GetClient(args, true)
	if err != nil {
		return nil, err
	}
	batchStatusResponse, batchStatusErr := mdataClient.Set(gtin, state, wait)

	if batchStatusErr != nil {
		return nil, batchStatusErr
	}

	// Query batch transaction status link
	status := commands.GetTransactionStatus(batchStatusResponse)

	return status, nil
}
