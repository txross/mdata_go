/**
 * Copyright 2017-2018 Intel Corporation
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

package handler

import (
	"fmt"
	"mdata_go/mdata_payload"
	"mdata_go/mdata_state"
	"strings"

	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

var logger *logging.Logger = logging.Get()

type MdHandler struct {
}

func (self *MdHandler) FamilyName() string {
	/* In Sawtooth, transactions are defined by an extensible system called "transaction families"
	 */
	return "mdata"
}

func (self *MdHandler) FamilyVersions() []string {
	// Versions allow you to correlate deployments among all the nodes in your  network. You want all the nodes using the same version
	return []string{"1.0"}
}

func (self *MdHandler) Namespaces() []string {
	/* Namespace prefix is six hex characters, or three bytes
	All data under a namespace prefix follows a consistent address and data encoding/serialization schem that is determined
	by the transaction family which defines the namespace
	*/
	return []string{mdata_state.Namespace}
}

//Signature before interfaces:
func (self *MdHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {

	// type TpProcessRequest struct {
	// 	Header               *transaction_pb2.TransactionHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// 	Payload              []byte                             `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	// 	Signature            string                             `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	// 	ContextId            string                             `protobuf:"bytes,4,opt,name=context_id,json=contextId,proto3" json:"context_id,omitempty"`
	// 	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	// 	XXX_unrecognized     []byte                             `json:"-"`
	// 	XXX_sizecache        int32                              `json:"-"`
	// }

	// type Context struct {
	// 	connection messaging.Connection
	// 	contextId  string
	// }

	// The master data organization is defined as the signer of the transaction, so we unpack
	// the transaction header to obtain the signer's public key, which will be
	// used as the organization's identity.
	header := request.GetHeader()         //sawtooth-sdk-go/protobuf/processor_pb2
	signer := header.GetSignerPublicKey() //sawtooth-sdk-go/protobuf/transaction_pb2

	// The payload is sent to the transaction processor as bytes (just as it
	// appears in the transaction constructed by the transactor).  We unpack
	// the payload into an MdPayload struct so we can access its fields.
	payload, err := mdata_payload.FromBytes(request.GetPayload())
	if err != nil {
		return err
	}

	// Context provides an abstract interface for getting and setting validator
	// state. All validator interactions by a handler should be through a Context
	// instance. Currently, the Context class is NOT thread-safe and Context classes
	// may not share the same messaging.Connection object.
	mdState := mdata_state.NewMdState(context)

	logger.Debugf("mdata txn %v: signer %v: payload: Action='%v', Gtin='%v', Material='%v'",
		request.GetSignature(), signer, payload.Gtin, payload.Attributes)

	switch payload.Action {
	case "create":
		err := validateCreate(mdState, payload.Gtin)
		if err != nil {
			return err
		}
		product := &mdata_state.Product{
			Gtin:       payload.Gtin,
			Attributes: mdata_state.deserializeAttributes(payload.Attributes),
			State:      "ACTIVE",
		}
		displayCreate(payload, signer)
		return mdState.SetProduct(payload.Gtin, product)
	case "delete":
		err := validateDelete(mdState, payload.Gtin)
		if err != nil {
			return err
		}
		return mdState.DeleteProduct(payload.Gtin)
	case "update":
		err := validateUpdate(mdState, payload.Gtin)
		if err != nil {
			return err
		}
		product, _ := mdState.GetProduct(payload.Gtin) //err is not needed here, as it is checked in the validateUpdate function
		product.Attributes = mdata_state.deserializeAttributes(payload.Attributes)
		product.State = "ACTIVE"
		displayUpdate(payload, signer, product)
		return mdState.SetProduct(payload.Gtin, product)
	case "deactivate":
		err := validateDeactivate(mdState, payload.Gtin)
		if err != nil {
			return err
		}
		product, _ := mdState.GetProduct(payload.Gtin) //err is not needed here, as it is checked in the validateDeactivate function
		product.State = "INACTIVE"
		displayDeactivate(payload, signer, product)
		return mdState.SetProduct(payload.Gtin, product)
	default:
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("Invalid Action : '%v'", payload.Action)}
	}
}

func validateCreate(mdState *mdata_state.MdState, gtin string) error {
	product, err := mdState.GetProduct(gtin)
	if err != nil {
		return err
	}
	if product != nil {
		return &processor.InvalidTransactionError{Msg: "Product already exists"}
	}

	return nil
}

func displayCreate(payload *mdata_payload.MdPayload, signer string) {
	s := fmt.Sprintf("+ Signer %s created product %s +", signer[:6], payload.Gtin)
	sLength := len(s)
	border := "+" + strings.Repeat("-", sLength-2) + "+"
	fmt.Println(border)
	fmt.Println(s)
	fmt.Println(border)
}

func validateUpdate(mdState *mdata_state.MdState, gtin string) error {
	product, err := mdState.GetProduct(gtin)
	if err != nil {
		return err
	}
	if product == nil {
		return &processor.InvalidTransactionError{Msg: "Update requires an existing product"}
	}
	return nil
}

func displayUpdate(payload *mdata_payload.MdPayload, signer string, product *mdata_state.Product) {
	s := fmt.Sprintf("+ Signer %s updated product %s with attributes %s", signer[:6], product.Gtin, product.Attributes)
	sLength := len(s)
	border := "+" + strings.Repeat("-", sLength-2) + "+"
	fmt.Println(border)
	fmt.Println(s)
	fmt.Println(border)
}

func validateDeactivate(mdState *mdata_state.MdState, gtin string) error {
	product, err := mdState.GetProduct(gtin)
	if err != nil {
		return err
	}
	if product == nil {
		return &processor.InvalidTransactionError{Msg: "Deactivate requires an existing product"}
	}
	return nil
}

func displayDeactivate(payload *mdata_payload.MdPayload, signer string, product *mdata_state.Product) {
	s := fmt.Sprintf("+ Signer %s updated product %s state to %s+", signer[:6], product.Gtin, product.State)
	sLength := len(s)
	border := "+" + strings.Repeat("-", sLength-2) + "+"
	fmt.Println(border)
	fmt.Println(s)
	fmt.Println(border)
}

func validateDelete(mdState *mdata_state.MdState, gtin string) error {
	product, err := mdState.GetProduct(gtin)
	if err != nil {
		return err
	}
	if product == nil {
		return &processor.InvalidTransactionError{Msg: "Delete requires an existing product"}
	}
	if product.State != "INACTIVE" {
		return &processor.InvalidTransactionError{Msg: "Delete requires an INACTIVE product. Please deactivate the product with `mdata deactivate <GTIN>`."}
	}
	return nil
}
