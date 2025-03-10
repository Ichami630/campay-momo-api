package helper

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// prompt user
func GetInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n') //detects when a user goes to new line
	input = strings.TrimSpace(input) //removes trailing and white spaces in string

	return input, err
}

// check status code
func StatusCode(status int) {
	if status != http.StatusOK {
		log.Fatalf("Unexpected status response error %d - %s", status, http.StatusText(status))
	}
}

// validate phone number
func IsDigitisOnly(s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
