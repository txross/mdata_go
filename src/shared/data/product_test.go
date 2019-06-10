package data

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var testGtin1 string = "11111111111111"
var testAttributesEmpty Attributes = Attributes{}
var testState string = "ACTIVE"
var testAttributesOne Attributes = Attributes{"uom": "cases"}
var testGtin2 string = "55555555555555"
var testAttributesMulti Attributes = Attributes{"uom": "lbs", "weight": "300"}

func TestSerializedAttributes(t *testing.T) {

	tests := map[string]struct {
		attr          Attributes
		outSerialized []byte
	}{
		"nilAttribute": {
			attr:          testAttributesEmpty,
			outSerialized: []byte(nil),
		},
		"oneAttribute": {
			attr:          testAttributesOne,
			outSerialized: []byte("uom=cases"),
		},
		"multiAttribute": {
			attr:          testAttributesMulti,
			outSerialized: []byte("uom=lbs,weight=300"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		serialized := test.attr.Serialize()
		assert.Equal(t, test.outSerialized, serialized)
	}
}

func TestDeserializedAttributes(t *testing.T) {

	tests := map[string]struct {
		in              []string
		outDeserialized Attributes
	}{
		"nilAttribute": {
			in:              []string{},
			outDeserialized: Attributes{},
		},
		"oneAttribute": {
			in:              []string{"uom=lbs"},
			outDeserialized: Attributes{"uom": "lbs"},
		},
		"multiAttribute": {
			in:              []string{"uom=lbs,weight=300"},
			outDeserialized: Attributes{"uom": "lbs", "weight": "300"},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		deserialized := DeserializeAttributes(test.in)
		assert.Equal(t, test.outDeserialized, deserialized)
	}
}

var testProduct Product = Product{
	Gtin:       testGtin1,
	Attributes: testAttributesEmpty,
	State:      testState,
}

var testProduct2 Product = Product{
	Gtin:       testGtin2,
	Attributes: testAttributesMulti,
	State:      testState,
}

var testProductSliceEmpty []*Product = []*Product{}

var testProductSliceOne []*Product = []*Product{&testProduct}

var testProductSliceMulti []*Product = []*Product{
	&testProduct,
	&testProduct2,
}

func TestSerializedProduct(t *testing.T) {

	tests := map[string]struct {
		in            []*Product
		outSerialized []byte
	}{
		"nilProduct": {
			in:            testProductSliceEmpty,
			outSerialized: []byte(nil),
		},
		"onProduct": {
			in:            testProductSliceOne,
			outSerialized: []byte(fmt.Sprintf("%v,%v,%v", testProduct.Gtin, string(testProduct.Attributes.Serialize()), testProduct.State)),
		},
		"multiProduct": {
			in: testProductSliceMulti,
			outSerialized: []byte(
				fmt.Sprintf("%v,%v,%v|%v,%v,%v",
					testProduct.Gtin,
					string(testProduct.Attributes.Serialize()),
					testProduct.State,
					testProduct2.Gtin,
					string(testProduct2.Attributes.Serialize()),
					testProduct2.State)),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		serialized := Serialize(test.in)
		assert.Equal(t, test.outSerialized, serialized)
	}
}

func TestDeserializedProduct(t *testing.T) {

	tests := map[string]struct {
		in              []byte
		outDeserialized map[string]*Product
		outErr          error
	}{
		"nilProduct": {
			in:              Serialize(testProductSliceEmpty),
			outDeserialized: map[string]*Product(nil),
			outErr:          errors.New("Malformed product"),
		},
		"onProduct": {
			in: Serialize(testProductSliceOne),
			outDeserialized: map[string]*Product{
				testGtin1: &testProduct,
			},
			outErr: nil,
		},
		"multiProduct": {
			in: Serialize(testProductSliceMulti),
			outDeserialized: map[string]*Product{
				testGtin1: &testProduct,
				testGtin2: &testProduct2,
			},
			outErr: nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		deserialized, err := Deserialize(test.in)
		assert.Equal(t, test.outDeserialized, deserialized)
		assert.Equal(t, reflect.TypeOf(test.outErr), reflect.TypeOf(err))
	}
}
