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
	"github.com/jessevdk/go-flags"
)

type SetState struct {
	Args struct {
		Gtin  string `positional-arg-name:"gtin" required:"true" description:"Identify the gtin of the product to set state"`
		State string `positional-arg-name:"state" required:"true" description:"Specify the state to set the <gtin>: ACTIVE, INACTIVE, DISCONTINUED" choice:"INACTIVE" choice:"ACTIVE" choice:"DISCONTINUED"`
	} `positional-args:"true"`
	Url     string `long:"url" description:"Specify URL of REST API"`
	Keyfile string `long:"keyfile" description:"Identify file containing user's private key"`
	Wait    uint   `long:"wait" description:"Set time, in seconds, to wait for transaction to commit"`
}

func (args *SetState) Name() string {
	return "set"
}

func (args *SetState) KeyfilePassed() string {
	return args.Keyfile
}

func (args *SetState) UrlPassed() string {
	return args.Url
}

func (args *SetState) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "SetState for a product", "Sends an mdata transaction to set state of <gtin> to <state>.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *SetState) Run() error {
	// Construct client
	gtin := args.Args.Gtin
	state := args.Args.State
	wait := args.Wait

	mdataClient, err := GetClient(args, true)
	if err != nil {
		return err
	}
	_, err = mdataClient.SetState(gtin, state, wait)
	return err
}
