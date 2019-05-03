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

package create

import (
	flags "github.com/jessevdk/go-flags"
	"mdata_go/src/mdata_client/client"
)

type Create struct {
	Args struct {
		Gtin string `positional-arg-name:"gtin" required:"true" description:"Identify the gtin of the product to create"`
	} `positional-args:"true"`
	Attributes map[string]string `long:"attributes" short:"a" required:"false" description:"Specify key:value pair to define product attributes"`
	Url        string            `long:"url" description:"Specify URL of REST API"`
	Keyfile    string            `long:"keyfile" description:"Identify file containing user's private key"`
	Wait       uint              `long:"wait" description:"Set time, in seconds, to wait for transaction to commit"`
}

func (args *Create) Name() string {
	return "create"
}

func (args *Create) KeyfilePassed() string {
	return args.Keyfile
}

func (args *Create) UrlPassed() string {
	return args.Url
}

func (args *Create) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Creates an product", "Sends an mdata transaction to create <gtin> with <attributes>.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *Create) Run() error {
	// Construct client
	gtin := args.Args.Gtin
	attributes := args.Attributes
	wait := args.Wait

	mdataClient, err := GetClient(args, true)
	if err != nil {
		return err
	}
	_, err = mdataClient.Create(gtin, attributes, wait)
	return err
}
