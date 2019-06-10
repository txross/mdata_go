package mdata_state

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	_data "github.com/tross-tyson/mdata_go/src/shared/data"
	"testing"
)

var testGtin string = "01234567891234"
var testAttributes _data.Attributes = _data.Attributes{"uom": "cases"}
var testSetNewAttributes _data.Attributes = _data.Attributes{"uom": "lbs", "weight": "300"}
var testState string = "ACTIVE"
var testGtinAddress string = makeAddress(testGtin)
var toDeleteGtin string = "555555555555"
var toDeleteGtinAddress string = makeAddress(toDeleteGtin)
var testProduct _data.Product = _data.Product{
	Gtin:       testGtin,
	Attributes: testAttributes,
	State:      testState,
}
var testSetNewProduct _data.Product = _data.Product{
	Gtin:       testGtin,
	Attributes: testSetNewAttributes,
	State:      testState,
}
var sampleError = errors.New("sample")

func TestGetProduct(t *testing.T) {

	tests := map[string]struct {
		gtin       string
		outProduct *_data.Product
		err        error
	}{
		"error": {
			gtin:       testGtin,
			outProduct: nil,
			err:        sampleError,
		},
		"emptyProduct": {
			gtin:       testGtin,
			outProduct: nil,
			err:        nil,
		},
		"existingProduct": {
			gtin:       testGtin,
			outProduct: &testProduct,
			err:        nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		testContext := &mockContext{}

		if name == "existingProduct" {
			returnState := make(map[string][]byte)
			testProductSlice := make([]*_data.Product, 1)
			testProductSlice[0] = &testProduct

			fmt.Printf("Existing Product %v\n", testProductSlice[0])
			seri := _data.Serialize(testProductSlice)
			fmt.Printf("Serialized Product %v\n", string(seri))
			deseri, _ := _data.Deserialize(seri)
			fmt.Printf("DESERIALIZED PRODUCTS: \n\t %v", deseri)

			fmt.Printf("\nProduct: \n\t%v", &desri[testGtin])

			returnState[testGtinAddress] = _data.Serialize(testProductSlice)

			testContext.On("GetState", []string{testGtinAddress}).Return(
				returnState,
				nil,
			)
		}
		if name == "emptyProduct" {
			testContext.On("GetState", []string{testGtinAddress}).Return(
				nil,
				nil,
			)
		}
		if name == "error" {
			testContext.On("GetState", []string{testGtinAddress}).Return(
				nil,
				sampleError,
			)
		}

		testState := &MdState{
			context:      testContext,
			addressCache: make(map[string][]byte),
		}

		product, err := testState.GetProduct(test.gtin)
		assert.Equal(t, test.outProduct, product)
		assert.Equal(t, test.err, err)

		testContext.AssertExpectations(t)

	}
}

func TestSetProduct(t *testing.T) {

	tests := map[string]struct {
		gtin      string
		inProduct *_data.Product
		err       error
	}{
		"newProduct": {
			gtin:      testGtin,
			inProduct: &testProduct,
			err:       nil,
		},
		"updateProductState": {
			gtin:      testGtin,
			inProduct: &testSetNewProduct,
			err:       nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		testContext := &mockContext{}
		testProductSlice := []*_data.Product{&testProduct}

		if name == "newProduct" {
			returnState := make(map[string][]byte)
			testContext.On("GetState", []string{testGtinAddress}).Return(
				returnState,
				nil,
			)

			data := _data.Serialize(testProductSlice)
			testContext.On("SetState", map[string][]byte{testGtinAddress: data}).Return(
				[]string{testGtinAddress},
				nil,
			)
		}

		if name == "updateProductState" {
			returnState := make(map[string][]byte)
			returnState[testGtinAddress] = _data.Serialize(testProductSlice)
			testContext.On("GetState", []string{testGtinAddress}).Return(
				returnState,
				nil,
			)

			data := _data.Serialize([]*_data.Product{&testSetNewProduct})
			testContext.On("SetState", map[string][]byte{testGtinAddress: data}).Return(
				[]string{testGtinAddress},
				nil,
			)

		}

		testState := &MdState{
			context:      testContext,
			addressCache: make(map[string][]byte),
		}

		err := testState.SetProduct(test.gtin, test.inProduct)
		assert.Equal(t, test.err, err)
		testContext.AssertExpectations(t)

	}
}

func TestDeleteProduct(t *testing.T) {

	tests := map[string]struct {
		gtin string
		err  error
	}{
		"productDoesNotExist": { //Delete state at testGtin address
			gtin: testGtin,
			err:  nil,
		},
		"storeProductsWithoutDeleted": { //If other products exist, just storeProducts at the state address of the deleted Gtin
			gtin: toDeleteGtin,
			err:  nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		testContext := &mockContext{}

		testProductSlice := make([]*_data.Product, 2)
		testProductSlice[0] = &testProduct

		testProduct2 := _data.Product{
			Gtin:       toDeleteGtin,
			Attributes: _data.Attributes{"uom": "cases"},
			State:      "INACTIVE",
		}
		testProductSlice[1] = &testProduct2

		if name == "productDoesNotExist" {
			returnState := make(map[string][]byte) // Return empty map
			testContext.On("GetState", []string{testGtinAddress}).Return(
				returnState,
				nil,
			)

			testContext.On("DeleteState", []string{testGtinAddress}).Return(
				nil,
				nil,
			).Once()
		}

		if name == "storeProductsWithoutDeleted" {
			returnState := make(map[string][]byte)
			returnState[toDeleteGtinAddress] = _data.Serialize(testProductSlice)
			testContext.On("GetState", []string{toDeleteGtinAddress}).Return(
				returnState,
				nil,
			)

			data := _data.Serialize([]*_data.Product{&testProduct})
			testContext.On("SetState", map[string][]byte{toDeleteGtinAddress: data}).Return(
				[]string{toDeleteGtinAddress},
				nil,
			).Once()

		}

		testState := &MdState{
			context:      testContext,
			addressCache: make(map[string][]byte),
		}

		err := testState.DeleteProduct(test.gtin)
		assert.Equal(t, test.err, err)
		testContext.AssertExpectations(t)

	}

}
