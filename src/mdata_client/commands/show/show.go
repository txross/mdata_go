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

package show

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/tross-tyson/mdata_go/src/mdata_client/client"
	"github.com/tross-tyson/mdata_go/src/shared/data"
	"os"
)

type Show struct {
	Args struct {
		Gtin string `positional-arg-name:"gtin" required:"true" description:"Identify the gtin of the product to create"`
	} `positional-args:"true"`
	Url string `long:"url" description:"Specify URL of REST API"`
}

func (args *Show) Name() string {
	return "show"
}

func (args *Show) KeyfilePassed() string {
	return ""
}

func (args *Show) UrlPassed() string {
	return args.Url
}

func (args *Show) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Displays the specified mdata attributes", "Shows the attribues of the product <gtin>.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *Show) Run() ([]byte, error) {
	//TODO: Check back here after mdataClient.Show() has been defined
	// Construct client
	gtin := args.Args.Gtin
	mdataClient, err := client.GetClient(args, false)
	if err != nil {
		return nil, err
	}
	products, err := mdataClient.Show(gtin)
	if err != nil {
		return nil, err
	}

	productMap, deserialize_err := data.Deserialize([]byte(products))

	if deserialize_err != nil {
		return nil, deserialize_err
	}

	response := productMap[gtin].GetJson()

	os.Stdout.Write(response)
	return response, nil
}
