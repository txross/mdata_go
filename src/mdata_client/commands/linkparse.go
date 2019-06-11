package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Response struct {
	Link string `json:"link"`
}

type BatchStatusData struct {
	Id                  string   `json:"id"`
	InvalidTransactions []string `json:"invalid_transactions"`
	Status              string   `json:"status"`
}

type BatchStatusResponse struct {
	Data []BatchStatusData `json:"data"`
	Link string            `json:"link"`
}

func GetTransactionStatus(txn string) string {
	status := getStatus(getLink(txn))
	return status
}

func getLink(txn string) string {

	var response Response
	err := json.Unmarshal([]byte(txn), &response)

	if err != nil {
		return err.Error()
	}

	return response.Link
}

func getStatus(link string) string {

	var validLinkPattern = regexp.MustCompile(`^http:\/\/.*$`)

	if !validLinkPattern.MatchString(link) {
		return errors.New("Malformed or missing link to batch transaction id").Error()
	}

	resp, err := http.Get(link)
	if err != nil {
		return err.Error()
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad Response from REST Api %v", resp.Status).Error()
	} else {
		return queryStatus(resp)
	}

}

func queryStatus(resp *http.Response) *BatchStatusResponse {

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	response := BatchStatusResponse{}
	err = json.Unmarshal(body, &response)

	if err != nil {
		return err.Error()
	}

	return response
}
