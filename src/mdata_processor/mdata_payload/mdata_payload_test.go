package mdata_payload

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"reflect"
	"testing"
)

var sampleError = processor.InvalidTransactionError{Msg: "Sample Error"}

var testPayloads = map[string]struct {
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
	"nullPayload": { //Null payload => Err
		in:         nil,
		outPayload: nil,
		outError:   &sampleError,
	},
	"missingGtinCreate": { //Missing GTIN => Err
		in:         []byte("create"),
		outPayload: nil,
		outError:   &sampleError,
	},
	"missingGtinUpdate": { //Missing GTIN => Err
		in:         []byte("update,uom=cases"),
		outPayload: nil,
		outError:   &sampleError,
	},
	"missingAction": { //Missing action => Err
		in:         []byte(",00012345600012,uom=cases"),
		outPayload: nil,
		outError:   &sampleError,
	},
	"validAttributesCreate": { //Create with valid Attributes => Ok
		in:         []byte("create,00012345600012,uom=cases"),
		outPayload: &MdPayload{Action: "create", Gtin: "00012345600012", Attributes: []string{"uom=lbs"}},
		outError:   nil,
	},
	"validAttributesUpdate": { //Update with Attributes => Ok
		in:         []byte("update,00012345600012,uom=lbs,weight=300"),
		outPayload: &MdPayload{Action: "update", Gtin: "00012345600012", Attributes: []string{"uom=lbs", "weight=300"}},
		outError:   nil,
	},
	"noAttributesUpdate": { // Update with len(Attributes) < 1  => Err
		in:         []byte("update,00012345600012"),
		outPayload: nil,
		outError:   &sampleError,
	},
	"invalidCharGtin": { //Invalid character '|' => Err
		in:         []byte("update,000123|45600012,uom=lbs,weight=300"),
		outPayload: nil,
		outError:   &sampleError,
	},
	"invalidCharAttr": { //Invalid character '|'  => Err
		in:         []byte("update,00012345600012,uom=lbs,weight=3|00"),
		outPayload: nil,
		outError:   &sampleError,
	},
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
