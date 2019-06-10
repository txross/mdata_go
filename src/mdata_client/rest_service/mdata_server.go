package rest_service

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tross-tyson/mdata_go/src/mdata_client/parser"
)

// type ProductState int

// const (
// 	INACTIVE     ProductState = 0
// 	ACTIVE       ProductState = 1
// 	DISCONTINUED ProductState = 2
// )

// func (ps ProductState) String() string {
// 	switch ps {
// 	case INACTIVE:
// 		return "INACTIVE"
// 	case ACTIVE:
// 		return "ACTIVE"
// 	case DISCONTINUED:
// 		return "DISCONTINUED"
// 	default:
// 		return "UNKNOWN"
// 	}
// }

// func convertToState(s string) (ProductState, error) {
// 	switch strings.ToUpper(s) {
// 	case "INACTIVE":
// 		return INACTIVE, nil
// 	case "ACTIVE":
// 		return ACTIVE, nil
// 	case "DISCONTINUED":
// 		return DISCONTINUED, nil
// 	default:
// 		return INACTIVE, fmt.Errorf("Unknown product state string: %s", s)

// 	}
// }

// type ProductStateUpdate struct {
// 	State string `json:"new_state" xml:"new_state" form:"new_state" query:"new_state"`
// }

// func (p Product) isValid() bool {
// 	gtin := p.Gtin
// 	pattern := regexp.MustCompile(`^\d{14}$`)
// 	if !pattern.MatchString(gtin) {
// 		return false
// 	}

// 	//attributes are all optional, so no need to validate
// 	_, err := convertToState(p.State)
// 	return err == nil
// }

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
	p := &data.Product{}

	//1 Get data
	if err := c.Bind(p); err != nil {
		return err
	}

	//2 Supply arguments to parser
	args := []string{
		"create",
		p.Gtin,
	}

	//3 Split attributes into arguments for the parser, append to args
	attributes := p.Attributes.Serialize()
	for _, key_value_pair := range strings.Split(string(attributes), ",") {
		key_value_pair = strings.Replace(key_value_pair, "=", ":")
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

	response := &CrudResponse{Status: status, Product: *p}

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

	p := &data.Product{}

	//1 Get data
	if err := c.Bind(p); err != nil {
		return err
	}

	//2 Supply arguments to parser
	args := []string{
		"update",
		p.Gtin,
	}

	//3 Split attributes into arguments for the parser, append to args
	attributes := p.Attributes.Serialize()
	for _, key_value_pair := range strings.Split(string(attributes), ",") {
		key_value_pair = strings.Replace(key_value_pair, "=", ":")
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

	response := &CrudResponse{Status: status, Product: *p}

	return c.JSON(http.StatusOK, response)
}

func updateProductState(c echo.Context) error {
	// Use this function to update state or attributes of existing product
	// An update of attributes will overwrite existing attributes of the product

	p := &data.Product{}

	//1 Get data
	if err := c.Bind(p); err != nil {
		return err
	}

	//2 Supply arguments to parser
	args := []string{
		"set",
		p.Gtin,
		p.State,
	}

	_, err := p.ParseArgs(args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", err)
	}

	status, cmd_err := runCmd(p.Command.Active.Name)
	if cmd_err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing request, %v", cmd_err)
	}

	response := &CrudResponse{Status: status, Product: *p}

	return c.JSON(http.StatusOK, response)
}

func Run(port int) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) //for now open to all origins

	p := parser.GetParser()

	e.GET("/products/:gtin", showProduct)                  // show specific product
	e.GET("/products", listproduct)                        // list all products
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
