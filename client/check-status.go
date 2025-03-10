package client

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"example.com/m/helper"
)

type reference struct {
	Status string `json:"status"`
}

func CheckTransationStatus(token, baseUrl string, refOptional ...string) string {
	var ref string
	if len(refOptional) > 0 && refOptional[0] != "" {
		// Use the provided reference
		ref = refOptional[0]
	} else {
		// Prompt the user for reference
		reader := bufio.NewReader(os.Stdin)
		inputRef, _ := helper.GetInput("Enter the payment Reference Number: ", reader)
		ref = inputRef
	}

	resp, err := http.NewRequest(http.MethodGet, baseUrl+"transaction/"+ref+"/", nil)

	//set the request headers
	resp.Header.Set("Authorization", "Token "+token)
	resp.Header.Set("Content-Type", "application/json")

	//sent the request using http.defaultclient
	client := &http.Client{}
	response, err := client.Do(resp)

	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
	defer response.Body.Close()

	//read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//display the status code and text if not ok
	helper.StatusCode(response.StatusCode)

	//parse json data into struct
	var status reference
	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Fatalf("Error parsing JSON Body %v", err)
	}
	//access the transaction reference
	return status.Status
}
