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

package commands

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"mdata_go/src/mdata_client/client"
	"strings"
)

type List struct {
	Url string `long:"url" description:"Specify URL of REST API"`
}

func (args *List) Name() string {
	return "list"
}

func (args *List) KeyfilePassed() string {
	return ""
}

func (args *List) UrlPassed() string {
	return args.Url
}

func (args *List) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Displays all mdata products", "Shows the attributes of all gtins in mdata state.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *List) Run() error {

	//TODO: Check back here after List() has been defined in mdataClient
	// Construct client
	mdataClient, err := client.GetClient(args, false)
	if err != nil {
		return err
	}
	products, err := mdataClient.List()
	if err != nil {
		return err
	}

	fmt.Printf("%-15v\t%-40v\t%-10v\t\n", "GTIN", "ATTRIBUTES", "STATE")

	for _, product := range products {
		for _, str := range strings.Split(string(product), "|") {
			parts := strings.Split(string(str), ",")
			gtin := parts[0]
			attrs := parts[1 : len(parts)-1]
			state := parts[len(parts)-1]

			fmt.Printf("%-v\t%-v\t%-v\t\n", gtin, attrs, state)
		}
	}
	return nil
}
