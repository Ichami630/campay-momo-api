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

	"example.com/m/helper"
)

type responseLink struct {
	Link string `json:"link"`
}

func PaymentLink(token string, baseUrl string) {
	//prompt user to enter the details for the payment
	reader := bufio.NewReader(os.Stdin)

	number, _ := helper.GetInput("Please enter your momo Number: ", reader)
	//validate the phone number
	if len(number) != 9 || !helper.IsDigitisOnly(number) {
		fmt.Println("Invalid number. It must be 9 digits long.")
		RequestPayment(token, baseUrl) //callback the Request payment function
	}

	amount, _ := helper.GetInput("Please enter the amount: ", reader)
	description, _ := helper.GetInput("Payment Description: ", reader)
	firstName, _ := helper.GetInput("First Name: ", reader)
	lastName, _ := helper.GetInput("Last Name: ", reader)
	email, _ := helper.GetInput("Email: ", reader)
	//encode body data
	postBody, _ := json.Marshal(map[string]string{
		"amount":               amount,
		"currency":             "XAF",
		"from":                 "237" + number,
		"description":          description,
		"first_name":           firstName,
		"last_name":            lastName,
		"email":                email,
		"external_reference":   "",
		"redirect_url":         "https://brandonichami.com",
		"failure_redirect_url": "https://mgsmarttrading.com",
		"payment_options":      "MOMO",
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.NewRequest(http.MethodPost, baseUrl+"get_payment_link/", responseBody)
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
	var payLink responseLink
	err = json.Unmarshal(body, &payLink)
	if err != nil {
		log.Fatalf("Error parsing JSON Body %v", err)
	}

	fmt.Println("Payment Link: ", payLink.Link)

}
