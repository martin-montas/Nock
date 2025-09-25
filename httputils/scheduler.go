package httputils

import (
	"fmt"
	"net/http"
	"os"
)

func CheckStatuscodeFromURL(u string) int {
	response, err := http.Get(u)
	if err != nil {
		fmt.Printf("Domain is unreachable %s\n", u)
		os.Exit(1)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return response.StatusCode
	}
	return response.StatusCode
}

func IsForbidden(statusCode int) bool {
	switch statusCode {
	case 404:
		return true
	case 500:
		return true
	default:
		return false

	}
}
