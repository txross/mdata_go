package mdata_payload

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"reflect"
	"testing"
)

var sampleError = processor.InvalidTransactionError{Msg: "Sample Error"}

var testPayloads = []struct {
	in         []byte
	outPayload *MdPayload
	outError   error
}{
	/* Test Cases
	1. Null payload => Err
	2. Missing GTIN => Err
	3. Missing action => Err
	4. Invalid Attributes (not in key=value pairs) => Err
	5. Valid Attributes => Ok
	6. Update with Attributes => Ok
	7. Update with len(Attributes) < 1  => Err
	8. Invalid character '|'
	*/
	//Input, expected return MdPayload, expected return Error
	{nil, nil, &sampleError},                                 //Null payload => Err
	{[]byte("create"), nil, &sampleError},                    //Missing GTIN => Err
	{[]byte("update,uom=cases"), nil, &sampleError},          //Missing GTIN => Err
	{[]byte(",00012345600012,uom=cases"), nil, &sampleError}, //Missing action => Err
	{[]byte("create,00012345600012,uom=cases"), &MdPayload{Action: "create", Gtin: "00012345600012", Attributes: []string{"uom=lbs"}}, nil},                        //Valid Attributes => Ok
	{[]byte("update,00012345600012,uom=lbs,weight=300"), &MdPayload{Action: "update", Gtin: "00012345600012", Attributes: []string{"uom=lbs", "weight=300"}}, nil}, //Update with Attributes => Ok
	{[]byte("update,00012345600012"), nil, &sampleError},                     // Update with len(Attributes) < 1  => Err
	{[]byte("update,000123|45600012,uom=lbs,weight=300"), nil, &sampleError}, //Invalid character '|'
	{[]byte("update,00012345600012,uom=lbs,weight=3|00"), nil, &sampleError}, //Invalid character '|'
}

func compareExpectedActualError(expectedErr error, actualError error) bool {
	return reflect.TypeOf(expectedErr) == reflect.TypeOf(actualError)
}

func compareStructs(expected, actual MdPayload) bool {
	expected_fields := reflect.TypeOf(expected)
	expected_values := reflect.ValueOf(expected)
	num_expected_fields := expected_fields.NumField()

	actual_fields := reflect.TypeOf(actual)
	actual_values := reflect.ValueOf(actual)
	num_actual_fields := actual_fields.NumField()

	if num_expected_fields != num_actual_fields {
		return false
	}

	for i := 0; i < num_expected_fields; i++ {
		expected_field := expected_fields.Field(i)
		expected_value := expected_values.Field(i)
		actual_field := actual_fields.Field(i)
		actual_value := actual_values.Field(i)

		if expected_field.Name != actual_field.Name {
			return false
		}

		switch expected_value.Kind() {
		case reflect.String:
			ev := expected_value.String()
			av := actual_value.String()
			if ev != av {
				return false
			}
		case reflect.Int:
			ev := expected_value.Int()
			av := actual_value.Int()
			if ev != av {
				return false
			}
		}

	}
	return true
}

func compareExpectedActualPayload(expectedPayload *MdPayload, actualPayload *MdPayload) bool {
	var areEqual bool
	if expectedPayload != nil {
		areEqual = compareStructs(*expectedPayload, *actualPayload)
	} else {
		areEqual = reflect.TypeOf(expectedPayload) == reflect.TypeOf(actualPayload)
	}
	return areEqual
}

func TestFromBytes(t *testing.T) {
	for _, tt := range testPayloads {
		payload, err := FromBytes(tt.in)
		if compareExpectedActualPayload(tt.outPayload, payload) != true || compareExpectedActualError(tt.outError, err) != true {
			t.Errorf("FromBytes(%v) => GOT %v, %v, WANT %v, %v", tt.in, payload, err, tt.outPayload, tt.outError)
		}
	}
}
