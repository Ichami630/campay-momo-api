package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/m/helper"
)

type responseReference struct {
	Reference string `json:"reference"`
}

func RequestPayment(token string, baseUrl string) {
	//prompt user to enter the details for the payment
	reader := bufio.NewReader(os.Stdin)

	number, _ := helper.GetInput("Please enter your momo Number: ", reader)
	//validate the phone number
	if len(number) != 9 || !helper.IsDigitisOnly(number) {
		fmt.Println("Invalid number. It must be 9 digits long.")
		RequestPayment(token, baseUrl) //callback the Request payment function
	}

	amount, _ := helper.GetInput("Please enter the amount: ", reader)
	//validate amount
	if !helper.IsDigitisOnly(amount) {
		fmt.Println("Invalid amount. Must be numeric")
		amount, _ = helper.GetInput("Please enter the amount: ", reader)
	}
	description, _ := helper.GetInput("Payment Description: ", reader)
	//encode body data
	postBody, _ := json.Marshal(map[string]string{
		"amount":             amount,
		"currency":           "XAF",
		"from":               "237" + number,
		"description":        description,
		"external_reference": "",
		"external_user":      "",
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.NewRequest(http.MethodPost, baseUrl+"collect/", responseBody)
	//set request headers
	resp.Header.Set("Authorization", "Token "+token)
	resp.Header.Set("Content-Type", "application/json")

	//sent the request via http.DefaultClient
	client := &http.Client{}
	response, err := client.Do(resp)

	if err != nil {
		log.Fatalf("Unable to make request %v", err)
	}
	//ensures the response is closed when function exits
	defer response.Body.Close()

	//get the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Unable to get response body %v", err)
	}

	//display the status code and text if not ok
	helper.StatusCode(response.StatusCode)

	//parse json data into struct
	var reference responseReference
	err = json.Unmarshal(body, &reference)
	if err != nil {
		log.Fatalf("Error parsing JSON Body %v", err)
	}

	fmt.Println("Transation initiated, waiting for confirmation...")
	maxAttemps := 6 //max retries before giving up
	var status string
	for i := 1; i <= maxAttemps; i++ {
		status = CheckTransationStatus(token, baseUrl, reference.Reference)

		if status == "PENDING" {
			fmt.Println("Transaction status PENDING... still waiting...")
			time.Sleep(5 * time.Second) //wait before next check
			continue
		}

		fmt.Println("Transaction ", status)
		break
	}
	//callback
	if status == "PENDING" {
		fmt.Println("PENDING Transaction, No Action, Reference Number: ", reference.Reference)
	}

}
