package utils

import (
	"net/http"
	"os"
)

func SetBiteshipHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("BITESHIP_API_KEY"))
}
