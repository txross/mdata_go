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

package mdata_state

import (
	"crypto/sha512"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
	_data "github.com/tross-tyson/mdata_go/src/shared/data"
)

type context interface {
	GetState([]string) (map[string][]byte, error)
	DeleteState([]string) ([]string, error)
	SetState(map[string][]byte) ([]string, error)
}

/* Namespace prefix is six hex characters, or three bytes
All data under a namespace prefix follows a consistent address and data encoding/serialization schem that is determined
by the transaction family which defines the namespace
*/
var Namespace = hexdigest("mdata")[:6]

// MdState handles addressing, serialization, deserialization,
// and holding an addressCache of data at the address.
type MdState struct {
	context      context
	addressCache map[string][]byte
}

func NewMdState(context *processor.Context) *MdState {
	return &MdState{
		context:      context,
		addressCache: make(map[string][]byte),
	}
}

// Define states to store
func (self *MdState) GetProduct(gtin string) (*_data.Product, error) {
	products, err := self.loadProducts(gtin)
	if err != nil {
		return nil, err
	}
	product, ok := products[gtin]
	if ok {
		return product, nil
	}
	return nil, nil
}

func (self *MdState) SetProduct(gtin string, product *_data.Product) error {
	products, err := self.loadProducts(gtin)
	if err != nil {
		return err
	}
	products[gtin] = product
	return self.storeProducts(gtin, products)
}

// DeleteProduct deletes the product from state, handling
// hash collisions.
func (self *MdState) DeleteProduct(gtin string) error {
	products, err := self.loadProducts(gtin)
	if err != nil {
		return err
	}
	delete(products, gtin)
	if len(products) > 0 {
		return self.storeProducts(gtin, products)
	} else {
		return self.deleteProducts(gtin)
	}
}

func (self *MdState) storeProducts(gtin string, products map[string]*_data.Product) error {
	address := makeAddress(gtin)

	var gtins []string

	//for each Gtin (key) in map[string]*_data.Product
	for gtin := range products {
		//append gtin to gtins slice of string
		gtins = append(gtins, gtin)
	}
	sort.Strings(gtins)

	var p []*_data.Product
	for _, gtin := range gtins {
		p = append(p, products[gtin])
	}

	data := _data.Serialize(p)

	self.addressCache[address] = data

	_, err := self.context.SetState(map[string][]byte{
		address: data,
	})
	return err
}

func (self *MdState) loadProducts(gtin string) (map[string]*_data.Product, error) {
	address := makeAddress(gtin)
	data, ok := self.addressCache[address]
	if ok {
		if self.addressCache[address] != nil {
			products, _ := _data.Deserialize(data)
			return products, nil
		}
		return make(map[string]*_data.Product), nil
	}
	results, err := self.context.GetState([]string{address})
	if err != nil {
		return nil, err
	}
	if len(string(results[address])) > 0 {
		self.addressCache[address] = results[address]
		products, _ := _data.Deserialize(data)
		return products, nil
	}
	self.addressCache[address] = nil
	products := make(map[string]*_data.Product)
	return products, nil
}

func (self *MdState) deleteProducts(gtin string) error {
	address := makeAddress(gtin)

	_, err := self.context.DeleteState([]string{address})
	return err
}

func makeAddress(gtin string) string {
	return Namespace + hexdigest(gtin)[:64]
}

func hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}
