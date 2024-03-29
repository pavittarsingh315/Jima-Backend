package utils

import (
	"NeraJima/configs"
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var accountID, authToken, fromNumber = configs.EnvTwilioIDKeyFrom()

var twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{AccountSid: accountID, Password: authToken})

func SendRegistrationText(code int, number string) {
	message := fmt.Sprintf("Here is your NeraJima verification code: %d. Code expires in 5 minutes!", code)

	params := &openapi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom(fromNumber)
	params.SetBody(message)

	res, err := twilioClient.ApiV2010.CreateMessage(params)
	_ = res
	_ = err
}

func SendPasswordResetText(code int, number string) {
	message := fmt.Sprintf("Here is your NeraJima password reset code: %d. Code expires in 5 minutes!", code)

	params := &openapi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom(fromNumber)
	params.SetBody(message)

	res, err := twilioClient.ApiV2010.CreateMessage(params)
	_ = res
	_ = err
}
