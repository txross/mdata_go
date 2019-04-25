package mdata_payload

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"reflect"
	"strings"
)

type MdPayload struct {
	Action string
	Gtin   string
	Mtrl   string
}

func (p MdPayload) CheckForInvaildChar() (bool, string) {
	values := reflect.ValueOf(p)
	num := values.NumField()

	for i := 0; i < num; i++ {
		v := values.Field(i)
		if v.Kind() == reflect.String {
			strv := v.String()
			if strings.Contains(strv, "|") {
				return true, strv
			}
		}
	}

	return false, ""
}

func FromBytes(payloadData []byte) (*MdPayload, error) {
	if payloadData == nil {
		return nil, &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}

	parts := strings.Split(string(payloadData), ",")
	if len(parts) < 2 {
		return nil, &processor.InvalidTransactionError{Msg: "Payload is malformed"}
	}

	payload := MdPayload{}
	payload.Action = parts[0]
	payload.Gtin = parts[1]

	if len(payload.Action) < 1 {
		return nil, &processor.InvalidTransactionError{Msg: "Action is required"}
	}

	if len(payload.Gtin) < 1 {
		return nil, &processor.InvalidTransactionError{Msg: "Gtin is required"}
	}

	if payload.Action == "create" || payload.Action == "update" {
		if len(parts) < 3 || len(parts[2]) < 1 {
			return nil, &processor.InvalidTransactionError{Msg: "Mtrl is required for create and update"}
		}
		payload.Mtrl = parts[2]
	}

	isInvalid, invalidString := payload.CheckForInvaildChar()
	if isInvalid {
		return nil, &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("Invalid Name (char '|' not allowed): '%v'", invalidString)}
	}

	return &payload, nil
}
