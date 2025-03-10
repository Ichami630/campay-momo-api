package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type paymentResponse struct {
	Reference string `json:"reference"`
	UssdCode  string `json:"ussd_code"`
	Operator  string `json:"operator"`
}

type statusResponse struct {
	Status json.RawMessage `json:"status"`
}

func camPay(token string) string {
	postBody, _ := json.Marshal(map[string]string{
		"amount":             "5",
		"currency":           "XAF",
		"from":               "237651742492",
		"description":        "Test",
		"external_reference": "",
		"external_user":      "",
	})

	responseBody := bytes.NewBuffer(postBody)

	//create a new http request
	resp, err := http.NewRequest(http.MethodPost, "https://demo.campay.net/api/collect/", responseBody)
	if err != nil {
		log.Fatalln(err)
	}

	//set request headers
	resp.Header.Set("Authorization", "Token "+token)
	resp.Header.Set("Content-Type", "application/json")

	//send the request via http.DefaultClient
	client := &http.Client{}
	response, err := client.Do(resp)

	//handle error
	if err != nil {
		log.Fatalf("An error occurred %v", err)
	}
	defer resp.Body.Close()

	//read the respponse
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//parse JSON Body into struct
	var paymentResp paymentResponse
	err = json.Unmarshal(body, &paymentResp)
	if err != nil {
		log.Fatalf("Error parsing JSON Body %v", err)
	}

	//access the transaction reference
	fmt.Println("Transation Reference:", paymentResp.Reference)
	return paymentResp.Reference
}

// function that checks the status of a transaction
func checkStatus(reference, token string) {
	for {
		req, err := http.NewRequest(http.MethodGet, "https://demo.campay.net/api/transaction/"+reference, nil)
		if err != nil {
			log.Fatalln(err)
		}

		//set the request headers
		req.Header.Set("Authorization", "Token "+token)
		req.Header.Set("Content-Type", "application/json")

		//sent the request using http.DefaultClient
		client := &http.Client{}
		response, err := client.Do(req)

		if err != nil {
			log.Fatalf("An error occured %v", err)
		}
		defer response.Body.Close()

		//read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		//unmarshal only the status field
		var statusResp statusResponse
		err = json.Unmarshal(body, &statusResp)

		// Convert json.RawMessage to a string
		var status string
		if err := json.Unmarshal(statusResp.Status, &status); err != nil {
			log.Fatalln("Error decoding status:", err)
		}

		// Print the extracted status
		fmt.Println("Transaction Status:", status)

		// Stop checking if the transaction is no longer pending
		if status != "PENDING" {
			break
		}

		// Wait for 5 seconds before checking again
		fmt.Println("Waiting for user to complete transaction...")
		time.Sleep(5 * time.Second)
	}
}
