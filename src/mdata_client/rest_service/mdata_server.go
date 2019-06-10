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

var p *flags.Parser = parser.GetParser()

type CrudResponse struct {
	Status  []byte       `json:"Status" sml:"Status" form:"Status" query:"Status"`
	Product data.Product `json:"Product" sml:"Product" form:"Product" query:"Product"`
}

func runCmd(cmd_name string) ([]byte, error) {
	for _, cmd := range parser.Commands() {
		if cmd.Name() == cmd_name {
			response, err := cmd.Run()
			return response, err
		}
	}
	return nil, nil
}

func showProduct(c echo.Context) error {
	//1. Get product id from REST param
	gtin := c.Param("gtin")

	//2 Supply arguments to parser
	args := []string{
		"show",
		gtin,
	}

	_, err := p.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", err)
	}

	response, cmd_err := runCmd(p.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", cmd_err)
	}

	return c.JSON(http.StatusOK, response)
}

func listProduct(c echo.Context) error {

	//2 Supply arguments to parser
	args := []string{
		"list",
	}

	_, err := p.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", err)
	}

	response, cmd_err := runCmd(p.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", cmd_err)
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

	_, err := p.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", err)
	}

	status, cmd_err := runCmd(p.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", cmd_err)
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

	_, err := p.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", err)
	}

	status, cmd_err := runCmd(p.Command.Active.Name)

	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", cmd_err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf(`{"Status": %v}`, status))
}

func updateProductAttributes(c echo.Context) error {
	// Use this function to update state or attributes of existing product
	// An update of attributes will overwrite existing attributes of the product

	product := &data.Product{}

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

	_, err := p.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", err)
	}

	status, cmd_err := runCmd(p.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", cmd_err)
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

	_, err := p.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", err)
	}

	status, cmd_err := runCmd(p.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", cmd_err)
	}

	response := &CrudResponse{Status: status, Product: *product}

	return c.JSON(http.StatusOK, response)
}

func Run(port int, parser *flags.Parser) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) //for now open to all origins

	e.GET("/products/:gtin", showProduct)                  // show specific product
	e.GET("/products", listProduct)                        // list all products
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
