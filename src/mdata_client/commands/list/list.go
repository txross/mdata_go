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

package list

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/tross-tyson/mdata_go/src/mdata_client/client"
	"github.com/tross-tyson/mdata_go/src/shared/data"
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

func (args *List) Run() (string, error) {

	//TODO: Check back here after List() has been defined in mdataClient
	// Construct client
	mdataClient, err := client.GetClient(args, false)
	if err != nil {
		return "", err
	}
	products, err := mdataClient.List()
	if err != nil {
		return "", err
	}

	productMap, _ := data.Deserialize([]byte(products))

	response := data.GetProductMapJson(productMap)

	return string(response), nil
}
