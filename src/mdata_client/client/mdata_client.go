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

package client

import (
	bytes2 "bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"github.com/tross-tyson/mdata_go/src/mdata_client/commands"  //mdata_client/commands
	"github.com/tross-tyson/mdata_go/src/mdata_client/constants" //mdata_client/constants
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"
)

var logger *logging.Logger = logging.Get()

func GetClient(args commands.Command, readFile bool) (MdataClient, error) {
	url := args.UrlPassed()
	if url == "" {
		url = constants.DEFAULT_URL
	}
	keyfile := ""
	if readFile {
		var err error
		keyfile, err = GetKeyfile(args.KeyfilePassed())
		if err != nil {
			return MdataClient{}, err
		}
	}
	return NewMdataClient(url, keyfile)
}

func GetKeyfile(keyfile string) (string, error) {
	if keyfile == "" {
		username, err := user.Current()
		if err != nil {
			return "", err
		}
		return path.Join(
			username.HomeDir, ".sawtooth", "keys", username.Username+".priv"), nil
	} else {
		return keyfile, nil
	}
}

type MdataClient struct {
	url    string
	signer *signing.Signer
}

type MdataClientAction struct {
	action string
	gtin   string
	wait   uint
	attrs  map[string]string
	state  string
}

func (c *MdataClientAction) serializePayload() string {
	//Convert map[string]string to []string{key=value key=value...} because that is what is expected by the processor payload
	attributes := []string{}
	for k, v := range c.attrs {
		attributes = append(attributes, fmt.Sprintf("%v=%v", k, v))
	}

	return fmt.Sprintf("%v,%v,%v,%v",
		c.action,
		c.gtin,
		strings.Join(attributes, ","),
		c.state)
}

func NewMdataClient(url string, keyfile string) (MdataClient, error) {

	var privateKey signing.PrivateKey
	if keyfile != "" {
		// Read private key file
		privateKeyStr, err := ioutil.ReadFile(keyfile)
		if err != nil {
			return MdataClient{},
				fmt.Errorf("Failed to read private key: %v", err)
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
	// Requires gtin, sets state to ACTIVE, attributes are optional
	gtin string, attrs map[string]string, wait uint) (string, error) {
	c := MdataClientAction{}
	c.action = constants.VERB_CREATE
	c.gtin = gtin
	c.wait = wait
	if len(attrs) > 0 {
		c.attrs = attrs
	} else {
		c.attrs = make(map[string]string)
	}
	c.state = ""
	return mdataClient.sendTransaction(c, wait)
}

func (mdataClient MdataClient) Update(
	// Requires gtin and attributes
	gtin string, attrs map[string]string, wait uint) (string, error) {
	c := MdataClientAction{}
	c.action = constants.VERB_UPDATE
	c.gtin = gtin
	c.wait = wait
	c.attrs = attrs
	c.state = ""
	return mdataClient.sendTransaction(c, wait)
}

func (mdataClient MdataClient) Delete(
	// Requires gtin
	gtin string, wait uint) (string, error) {
	c := MdataClientAction{}
	c.action = constants.VERB_DELETE
	c.gtin = gtin
	c.wait = wait
	c.attrs = make(map[string]string)
	c.state = ""
	return mdataClient.sendTransaction(c, wait)
}

func (mdataClient MdataClient) Set(
	// Requires gtin and state to change to
	gtin string, state string, wait uint) (string, error) {
	c := MdataClientAction{}
	c.action = constants.VERB_SET_STATE
	c.gtin = gtin
	c.wait = wait
	c.attrs = make(map[string]string)
	c.state = state
	return mdataClient.sendTransaction(c, wait)
}

func (mdataClient MdataClient) List() ([]byte, error) {

	// API to call
	apiSuffix := fmt.Sprintf("%s?address=%s",
		constants.STATE_API, mdataClient.getPrefix())
	response, err := mdataClient.sendRequest(apiSuffix, []byte{}, "", "")
	if err != nil {
		return []byte{}, err
	}

	var toReturn bytes2.Buffer

	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return nil,
			fmt.Errorf("Error reading response: %v", err)
	}
	encodedEntries := responseMap["data"].([]interface{})

	for index, entry := range encodedEntries {
		entryData, ok := entry.(map[interface{}]interface{})
		if !ok {
			return nil,
				errors.New("Error reading entry data")
		}

		stringData, ok := entryData["data"].(string)
		if !ok {
			return nil,
				errors.New("Error reading string data")
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(stringData)
		if err != nil {
			return nil,
				fmt.Errorf("Error decoding: %v", err)
		}

		toReturn.Write(decodedBytes)
		if index < len(encodedEntries)-1 {
			toReturn.WriteString("|")
		}
	}

	return toReturn.Bytes(), nil
}

func (mdataClient MdataClient) Show(gtin string) (string, error) {

	apiSuffix := fmt.Sprintf("%s/%s", constants.STATE_API, mdataClient.getAddress(gtin))
	response, err := mdataClient.sendRequest(apiSuffix, []byte{}, "", gtin)
	if err != nil {
		return "", err
	}
	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}
	data, ok := responseMap["data"].(string)
	if !ok {
		return "", errors.New("Error reading as string")
	}
	responseData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("Error decoding response: %v", err)
	}

	strData := string(responseData)

	return fmt.Sprintf("%v", strData), nil
}

func (mdataClient MdataClient) getStatus(
	batchId string, wait uint) (string, error) {

	// API to call
	apiSuffix := fmt.Sprintf("%s?id=%s&wait=%d",
		constants.BATCH_STATUS_API, batchId, wait)
	response, err := mdataClient.sendRequest(apiSuffix, []byte{}, "", "")
	if err != nil {
		return "", err
	}

	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
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
		return "", fmt.Errorf("Failed to connect to REST API: %v", err)
	}
	if response.StatusCode == 404 {
		logger.Debug(fmt.Sprintf("%v", response))
		return "", fmt.Errorf("No such product: %s", gtin)
	} else if response.StatusCode >= 400 {
		return "", fmt.Errorf("Error %d: %s", response.StatusCode, response.Status)
	}
	defer response.Body.Close()
	reponseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}
	return string(reponseBody), nil
}

func (mdataClient MdataClient) sendTransaction(c MdataClientAction, wait uint) (string, error) {
	payload := c.serializePayload()
	gtin := c.gtin
	// construct the address
	address := mdataClient.getAddress(gtin)

	// Construct TransactionHeader
	rawTransactionHeader := transaction_pb2.TransactionHeader{
		SignerPublicKey:  mdataClient.signer.GetPublicKey().AsHex(),
		FamilyName:       constants.FAMILY_NAME,
		FamilyVersion:    constants.FAMILY_VERSION,
		Dependencies:     []string{}, // empty dependency list
		Nonce:            strconv.Itoa(rand.Int()),
		BatcherPublicKey: mdataClient.signer.GetPublicKey().AsHex(),
		Inputs:           []string{address},
		Outputs:          []string{address},
		PayloadSha512:    Sha512HashValue(payload),
	}
	transactionHeader, err := proto.Marshal(&rawTransactionHeader)
	if err != nil {
		return "", fmt.Errorf("Unable to serialize transaction header: %v", err)
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
		return "", fmt.Errorf("Unable to construct batch list: %v", err)
	}
	batchId := rawBatchList.Batches[0].HeaderSignature
	batchList, err := proto.Marshal(&rawBatchList)
	if err != nil {
		return "", fmt.Errorf("Unable to serialize batch list: %v", err)
	}

	if wait > 0 {
		waitTime := uint(0)
		startTime := time.Now()
		response, err := mdataClient.sendRequest(
			constants.BATCH_SUBMIT_API, batchList, constants.CONTENT_TYPE_OCTET_STREAM, gtin)
		if err != nil {
			return "", err
		}
		for waitTime < wait {
			status, err := mdataClient.getStatus(batchId, wait-waitTime)
			if err != nil {
				return "", err
			}
			waitTime = uint(time.Since(startTime))
			if status != "PENDING" {
				return response, nil
			}
		}
		return response, nil
	}

	return mdataClient.sendRequest(
		constants.BATCH_SUBMIT_API, batchList, constants.CONTENT_TYPE_OCTET_STREAM, gtin)
}

func (mdataClient MdataClient) getPrefix() string {
	return Sha512HashValue(constants.FAMILY_NAME)[:constants.FAMILY_NAMESPACE_ADDRESS_LENGTH]
}

func (mdataClient MdataClient) getAddress(gtin string) string {
	prefix := mdataClient.getPrefix()
	productAddress := Sha512HashValue(gtin)[:constants.FAMILY_VERB_ADDRESS_LENGTH]
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
		return batch_pb2.BatchList{}, fmt.Errorf("Unable to serialize batch header: %v", err)
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

func Sha512HashValue(value string) string {
	hashHandler := sha512.New()
	hashHandler.Write([]byte(value))
	return strings.ToLower(hex.EncodeToString(hashHandler.Sum(nil)))
}
