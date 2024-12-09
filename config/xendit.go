package config

import (
	"os"

	"github.com/xendit/xendit-go/v6"
)

var XenditClient *xendit.APIClient

func InitXendit() {
	XenditClient = xendit.NewClient(os.Getenv("XENDIT_API_KEY"))
}
