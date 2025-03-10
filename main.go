package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	// "strconv"

	"example.com/m/client"
	"example.com/m/helper"
	"github.com/joho/godotenv"
)

// function to allow users choose an option
func promptOption() {
	reader := bufio.NewReader(os.Stdin)
	opt, _ := helper.GetInput("Please Choose an option (1- Make a Payment 2- Request Payment Link 3- Check Transaction Status 4- Airtime Transfer or type exit to quit):", reader)

	//handle exit case
	if opt == "exit" {
		fmt.Println("Exiting... Goodbye 👋")
		os.Exit(0)
	}
	//env files
	token, url := loadEnv()
	//convert the chosen option from string to int
	// option,err = strconv.ParseInt(opt,8,64)
	// if err != nil {
	// 	fmt.Println("Make you sure enter a valid interger option")
	// }
	switch opt {
	case "1":
		client.RequestPayment(token, url)
		promptOption()
	case "2":
		client.PaymentLink(token, url)
		promptOption()
	case "3":
		status := client.CheckTransationStatus(token, url)
		fmt.Println("Transaction Status:", status)
		promptOption()
	case "4":
		fmt.Println("This feature is not supported")
		promptOption()
	default:
		fmt.Println("Invalid option, please try again...")
		promptOption()
	}
}

func main() {

	promptOption()

}

// load env file
func loadEnv() (string, string) {
	//load the env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the token from environment variables
	authToken := os.Getenv("AUTH_TOKEN")
	baseUrl := os.Getenv("BASE_URL")
	if authToken == "" || baseUrl == "" {
		log.Fatal("Authorization token or URL is missing from the .env file")
	}

	return authToken, baseUrl
}
