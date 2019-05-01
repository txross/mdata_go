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
	bytes2 "bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MdataClient struct {
	url    string
	signer *signing.Signer
}

func NewMdataClient(url string, keyfile string) (MdataClient, error) {

	var privateKey signing.PrivateKey
	if keyfile != "" {
		// Read private key file
		privateKeyStr, err := ioutil.ReadFile(keyfile)
		if err != nil {
			return MdataClient{},
				errors.New(fmt.Sprintf("Failed to read private key: %v", err))
		}
		// Get private key object
		privateKey = signing.NewSecp256k1PrivateKey(privateKeyStr)
	} else {
		privateKey = signing.NewSecp256k1Context().NewRandomPrivateKey()
	}
	cryptoFactory := signing.NewCryptoFactory(signing.NewSecp256k1Context())
	signer := cryptoFactory.NewSigner(privateKey)
	return MdataClient{url, signer}, nil
}

func (mdataClient MdataClient) Create(
	gtin string, attrs map[string]string, wait uint) (string, error) {
	return mdataClient.sendTransaction(VERB_CREATE, gtin, attrs, wait)
}

func (mdataClient MdataClient) Update(
	gtin string, attrs map[string]string, wait uint) (string, error) {
	return mdataClient.sendTransaction(VERB_UPDATE, gtin, attrs, wait)
}

func (mdataClient MdataClient) Delete(
	gtin string, wait uint) (string, error) {
	return mdataClient.sendTransaction(VERB_DELETE, gtin, wait)
}

func (mdataClient MdataClient) SetState(
	gtin string, state string, wait uint) (string, error) {
	return mdataClient.sendTransaction(VERB_SET_STATE, gtin, state, wait)
}

func (mdataClient MdataClient) List() ([]byte, error) {

	// API to call
	apiSuffix := fmt.Sprintf("%s?address=%s",
		STATE_API, mdataClient.getPrefix())
	response, err := mdataClient.sendRequest(apiSuffix, []byte{}, "", "")
	if err != nil {
		return []byte{}, err
	}

	var toReturn []byte

	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return []map[interface{}]interface{}{},
			errors.New(fmt.Sprintf("Error reading response: %v", err))
	}
	encodedEntries := responseMap["data"].([]interface{})

	for _, entry := range encodedEntries {
		entryData, ok := entry.(map[interface{}]interface{})
		if !ok {
			return []byte{},
				errors.New("Error reading entry data")
		}

		stringData, ok := entryData["data"].(string)
		if !ok {
			return []byte{},
				errors.New("Error reading string data")
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(stringData)
		if err != nil {
			return []byte{},
				errors.New(fmt.Sprint("Error decoding: %v", err))
		}

		toReturn = append(toReturn, decodedBytes)
	}
	return toReturn, nil
}

func (mdataClient MdataClient) Show(gtin string) (string, error) {

	//TODO: Compare to python return of data from list

	apiSuffix := fmt.Sprintf("%s/%s", STATE_API, mdataClient.getAddress(gtin))
	response, err := mdataClient.sendRequest(apiSuffix, []byte{}, "", gtin)
	if err != nil {
		return "", err
	}
	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return "", errors.New(fmt.Sprint("Error reading response: %v", err))
	}
	data, ok := responseMap["data"].(string)
	if !ok {
		return "", errors.New("Error reading as string")
	}
	responseData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New(fmt.Sprint("Error decoding response: %v", err))
	}
	responseFinal := make(map[interface{}]interface{})
	err = cbor.Loads(responseData, &responseFinal)
	if err != nil {
		return "", errors.New(fmt.Sprint("Error binary decoding: %v", err))
	}
	return fmt.Sprintf("%v", responseFinal[gtin]), nil
}

func (mdataClient MdataClient) getStatus(
	batchId string, wait uint) (string, error) {

	// API to call
	apiSuffix := fmt.Sprintf("%s?id=%s&wait=%d",
		BATCH_STATUS_API, batchId, wait)
	response, err := mdataClient.sendRequest(apiSuffix, []byte{}, "", "")
	if err != nil {
		return "", err
	}

	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading response: %v", err))
	}
	entry :=
		responseMap["data"].([]interface{})[0].(map[interface{}]interface{})
	return fmt.Sprint(entry["status"]), nil
}

func (mdataClient MdataClient) sendRequest(
	apiSuffix string,
	data []byte,
	contentType string,
	gtin string) (string, error) {

	// Construct URL
	var url string
	if strings.HasPrefix(mdataClient.url, "http://") {
		url = fmt.Sprintf("%s/%s", mdataClient.url, apiSuffix)
	} else {
		url = fmt.Sprintf("http://%s/%s", mdataClient.url, apiSuffix)
	}

	// Send request to validator URL
	var response *http.Response
	var err error
	if len(data) > 0 {
		response, err = http.Post(url, contentType, bytes2.NewBuffer(data))
	} else {
		response, err = http.Get(url)
	}
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Failed to connect to REST API: %v", err))
	}
	if response.StatusCode == 404 {
		logger.Debug(fmt.Sprintf("%v", response))
		return "", errors.New(fmt.Sprintf("No such product: %s", gtin))
	} else if response.StatusCode >= 400 {
		return "", errors.New(
			fmt.Sprintf("Error %d: %s", response.StatusCode, response.Status))
	}
	defer response.Body.Close()
	reponseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading response: %v", err))
	}
	return string(reponseBody), nil
}

func (mdataClient MdataClient) sendTransaction(
	verb string, gtin string, attrs map[string]string, wait uint) (string, error) {

	// construct the payload information in CBOR format
	payloadData := make(map[string]interface{})
	payloadData["Verb"] = verb
	payloadData["Gtin"] = gtin
	payloadData["Attributes"] = attrs
	payload, err := cbor.Dumps(payloadData)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to construct CBOR: %v", err))
	}

	// construct the address
	address := mdataClient.getAddress(gtin)

	// Construct TransactionHeader
	rawTransactionHeader := transaction_pb2.TransactionHeader{
		SignerPublicKey:  mdataClient.signer.GetPublicKey().AsHex(),
		FamilyName:       FAMILY_NAME,
		FamilyVersion:    FAMILY_VERSION,
		Dependencies:     []string{}, // empty dependency list
		Nonce:            strconv.Itoa(rand.Int()),
		BatcherPublicKey: mdataClient.signer.GetPublicKey().AsHex(),
		Inputs:           []string{address},
		Outputs:          []string{address},
		PayloadSha512:    Sha512HashValue(string(payload)),
	}
	transactionHeader, err := proto.Marshal(&rawTransactionHeader)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Unable to serialize transaction header: %v", err))
	}

	// Signature of TransactionHeader
	transactionHeaderSignature := hex.EncodeToString(
		mdataClient.signer.Sign(transactionHeader))

	// Construct Transaction
	transaction := transaction_pb2.Transaction{
		Header:          transactionHeader,
		HeaderSignature: transactionHeaderSignature,
		Payload:         []byte(payload),
	}

	// Get BatchList
	rawBatchList, err := mdataClient.createBatchList(
		[]*transaction_pb2.Transaction{&transaction})
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Unable to construct batch list: %v", err))
	}
	batchId := rawBatchList.Batches[0].HeaderSignature
	batchList, err := proto.Marshal(&rawBatchList)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Unable to serialize batch list: %v", err))
	}

	if wait > 0 {
		waitTime := uint(0)
		startTime := time.Now()
		response, err := mdataClient.sendRequest(
			BATCH_SUBMIT_API, batchList, CONTENT_TYPE_OCTET_STREAM, gtin)
		if err != nil {
			return "", err
		}
		for waitTime < wait {
			status, err := mdataClient.getStatus(batchId, wait-waitTime)
			if err != nil {
				return "", err
			}
			waitTime = uint(time.Now().Sub(startTime))
			if status != "PENDING" {
				return response, nil
			}
		}
		return response, nil
	}

	return mdataClient.sendRequest(
		BATCH_SUBMIT_API, batchList, CONTENT_TYPE_OCTET_STREAM, gtin)
}

func (mdataClient MdataClient) getPrefix() string {
	return Sha512HashValue(FAMILY_NAME)[:FAMILY_NAMESPACE_ADDRESS_LENGTH]
}

func (mdataClient MdataClient) getAddress(gtin string) string {
	prefix := mdataClient.getPrefix()
	productAddress := Sha512HashValue(gtin)[FAMILY_VERB_ADDRESS_LENGTH:]
	return prefix + productAddress
}

func (mdataClient MdataClient) createBatchList(
	transactions []*transaction_pb2.Transaction) (batch_pb2.BatchList, error) {

	// Get list of TransactionHeader signatures
	transactionSignatures := []string{}
	for _, transaction := range transactions {
		transactionSignatures =
			append(transactionSignatures, transaction.HeaderSignature)
	}

	// Construct BatchHeader
	rawBatchHeader := batch_pb2.BatchHeader{
		SignerPublicKey: mdataClient.signer.GetPublicKey().AsHex(),
		TransactionIds:  transactionSignatures,
	}
	batchHeader, err := proto.Marshal(&rawBatchHeader)
	if err != nil {
		return batch_pb2.BatchList{}, errors.New(
			fmt.Sprintf("Unable to serialize batch header: %v", err))
	}

	// Signature of BatchHeader
	batchHeaderSignature := hex.EncodeToString(
		mdataClient.signer.Sign(batchHeader))

	// Construct Batch
	batch := batch_pb2.Batch{
		Header:          batchHeader,
		Transactions:    transactions,
		HeaderSignature: batchHeaderSignature,
	}

	// Construct BatchList
	return batch_pb2.BatchList{
		Batches: []*batch_pb2.Batch{&batch},
	}, nil
}
