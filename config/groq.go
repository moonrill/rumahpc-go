package config

import (
	"net/http"
	"os"

	"github.com/magicx-ai/groq-go/groq"
)

var GroqClient groq.Client

func InitGroq() {
	GroqClient = groq.NewClient(os.Getenv("GROQ_API_KEY"), &http.Client{})
}
