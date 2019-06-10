package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Attributes map[string]interface{}

func (self Attributes) Serialize() []byte {
	var b bytes.Buffer
	var i int = 0
	for k, v := range self {
		b.WriteString(fmt.Sprintf("%v=%v", k, v))
		i += 1
		if i < len(self) {
			b.WriteString(",")
		}
	}
	return b.Bytes()
}

func DeserializeAttributes(a []string) Attributes {
	A := Attributes{}
	for _, str := range a {
		if str != "" {
			parts := strings.Split(str, "=")
			k, v := parts[0], parts[1]
			A[k] = v
		}
	}

	return A
}

func Deserialize(data []byte) (map[string]*Product, error) {
	products := make(map[string]*Product)
	for _, str := range strings.Split(string(data), "|") {
		parts := strings.Split(string(str), ",")
		if len(parts) < 3 { //Product must have at least three serialized attributes (even if Product.Attributes is empty)
			return nil, errors.New(fmt.Sprintf("Malformed product data: '%v'", string(data)))
		}
		attrs := parts[1 : len(parts)-1]

		product := &Product{
			Gtin:       parts[0],
			Attributes: DeserializeAttributes(attrs),
			State:      parts[len(parts)-1],
		}
		products[parts[0]] = product
	}
	return products, nil
}

func Serialize(products []*Product) []byte {
	var buffer bytes.Buffer
	for i, product := range products {
		//00001234567890,uom=cases,weight=200,ACTIVE|
		buffer.WriteString(product.Gtin)
		buffer.WriteString(",")
		buffer.WriteString(string(product.Attributes.Serialize()))
		buffer.WriteString(",")
		buffer.WriteString(product.State)
		if i+1 != len(products) {
			buffer.WriteString("|")
		}
	}
	return buffer.Bytes()
}

type Product struct {
	Gtin       string     `json:"gtin" sml:"gtin" form:"gtin" query:"gtin"`
	Attributes Attributes `json:"attributes" xml:"attributes" form:"attributes" query:"attributes"`
	State      string     `json:"state" xml:"state" form:"state" query:"state"`
}

func (p *Product) GetJson() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("Error marshalling product json, %v", err)
		return nil
	}
	return b
}

func GetProductMapJson(productMap map[string]*Product) []byte {
	b, err := json.Marshal(productMap)
	if err != nil {
		fmt.Printf("Error marshalling product json, %v", err)
		return nil
	}
	return b
}
