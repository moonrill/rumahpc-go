package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/moonrill/rumahpc-api/templates"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func GenerateOTP(length int) string {

	var digits = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var otp string

	for i := 0; i < length; i++ {
		otp += strconv.Itoa(digits[rand.Intn(len(digits))])
	}

	return otp
}

func SendOTPEmail(email, otp string) error {
	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	to := []string{email}

	headers := map[string]string{
		"From":         os.Getenv("SMTP_EMAIL"),
		"To":           email,
		"Subject":      "OTP Verification",
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var headerString string
	for key, value := range headers {
		headerString += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	formattedHtml := fmt.Sprintf(templates.OTP, otp, time.Now().Year())

	message := headerString + "\r\n\r\n" + formattedHtml

	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("SMTP_EMAIL"), to, []byte(message))
	return err
}

func SendOTPSMS(phoneNumber, otp string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(fmt.Sprintf("Your OTP is: %s", otp))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
		return err
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

	return nil
}
