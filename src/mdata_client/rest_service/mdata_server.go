package rest_service

import (
	"fmt"
	"net/http"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tross-tyson/mdata_go/src/mdata_client/parser"
	"github.com/tross-tyson/mdata_go/src/shared/data"
)

var RestServiceParser *flags.Parser = parser.GetParser(parser.Commands())

type CrudResponse struct {
	Status  string       `json:"Status" sml:"Status" form:"Status" query:"Status"`
	Product data.Product `json:"Product" sml:"Product" form:"Product" query:"Product"`
}

func runCmd(cmd_name string) (string, error) {
	for _, cmd := range parser.Commands() {
		if cmd.Name() == cmd_name {
			response, err := cmd.Run()
			return response, err
		}
	}
	return "", fmt.Errorf("Command active name not found %v", cmd_name)
}

func listProduct(c echo.Context) error {

	//2 Supply arguments to parser
	args := []string{
		"list",
	}

	_, err := RestServiceParser.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing arguments %v, %v", args, err)
	}

	response, cmd_err := runCmd(RestServiceParser.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error executing command %v, %v", RestServiceParser.Command.Active.Name, cmd_err)
	}

	return c.JSON(http.StatusOK, response)
}

func showProduct(c echo.Context) error {
	//1. Get product id from REST param
	gtin := c.Param("gtin")

	//2 Supply arguments to parser
	args := []string{
		"show",
		gtin,
	}

	_, err := RestServiceParser.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing arguments %v, %v", args, err)
	}

	response, cmd_err := runCmd(RestServiceParser.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error executing command %v, %v", RestServiceParser.Command.Active.Name, cmd_err)
	}

	return c.JSON(http.StatusOK, response)
}

func createProduct(c echo.Context) error {
	product := &data.Product{}

	//1 Get data
	if err := c.Bind(product); err != nil {
		return err
	}

	//2 Supply arguments to parser
	args := []string{
		"create",
		product.Gtin,
	}

	//3 Split attributes into arguments for the parser, append to args
	attributes := product.Attributes.Serialize()
	for _, key_value_pair := range strings.Split(string(attributes), ",") {
		key_value_pair = strings.Replace(key_value_pair, "=", ":", 1)
		args = append(args, "-a", key_value_pair)
	}

	_, err := RestServiceParser.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing arguments %v, %v", args, err)
	}

	status, cmd_err := runCmd(RestServiceParser.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error executing command %v, %v", RestServiceParser.Command.Active.Name, cmd_err)
	}

	response := &CrudResponse{Status: status, Product: *product}

	return c.JSON(http.StatusOK, response)
}

func deleteProduct(c echo.Context) error {
	// Use this function to delete an existing product
	// Product must be in INACTIVE state to delete

	//1 Get params
	gtin := c.Param("gtin")

	//2 Supply arguments to parser
	args := []string{
		"delete",
		gtin,
	}

	_, err := RestServiceParser.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing arguments %v, %v", args, err)
	}

	status, cmd_err := runCmd(RestServiceParser.Command.Active.Name)

	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error executing command %v, %v", RestServiceParser.Command.Active.Name, cmd_err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf(`{"Status": %v}`, status))
}

func updateProductAttributes(c echo.Context) error {
	// Use this function to update state or attributes of existing product
	// An update of attributes will overwrite existing attributes of the product

	product := &data.Product{}

	//i.e.
	/*

		SAMPLE EXPECTED HTTP REQUEST JSON
		request_data : {
			Product: <gtin>,
			Atributes: {
				<key1>: <value1>,
				<key2>: <value2>,
				...
				<keyn>: <valuen>
			},
			State: <state>
		}

	*/

	//1 Get data
	if err := c.Bind(product); err != nil {
		return err
	}

	//2 Supply arguments to parser
	args := []string{
		"update",
		product.Gtin,
	}

	//3 Split attributes into arguments for the parser, append to args
	attributes := product.Attributes.Serialize()
	for _, key_value_pair := range strings.Split(string(attributes), ",") {
		key_value_pair = strings.Replace(key_value_pair, "=", ":", 1)
		args = append(args, "-a", key_value_pair)
	}

	_, err := RestServiceParser.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing arguments %v, %v", args, err)
	}

	status, cmd_err := runCmd(RestServiceParser.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error executing command %v, %v", RestServiceParser.Command.Active.Name, cmd_err)
	}

	response := &CrudResponse{Status: status, Product: *product}

	return c.JSON(http.StatusOK, response)
}

func updateProductState(c echo.Context) error {
	// Use this function to update state or attributes of existing product
	// An update of attributes will overwrite existing attributes of the product

	product := &data.Product{}

	//1 Get data
	if err := c.Bind(product); err != nil {
		return err
	}

	//2 Supply arguments to parser
	args := []string{
		"set",
		product.Gtin,
		product.State,
	}

	_, err := RestServiceParser.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing arguments %v, %v", args, err)
	}

	status, cmd_err := runCmd(RestServiceParser.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error executing command %v, %v", RestServiceParser.Command.Active.Name, cmd_err)
	}

	response := &CrudResponse{Status: status, Product: *product}

	return c.JSON(http.StatusOK, response)
}

func Run(port uint) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) //for now open to all origins

	e.GET("/products", listProduct)       // list all products
	e.GET("/products/:gtin", showProduct) // show specific product

	e.POST("/products", createProduct)                     // create new product
	e.PUT("/products/:gtin/attr", updateProductAttributes) // update existing product attributes or state
	e.PUT("/products/:gtin/state", updateProductState)     // update existing product attributes or state
	e.DELETE("/products/:gtin", deleteProduct)             // delete existing inactive product

	if port != 0 {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
	} else {
		e.Logger.Fatal(e.Start(":8888"))
	}

}
