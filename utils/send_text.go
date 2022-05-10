package utils

import (
	"NeraJima/configs"
	"encoding/json"
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
	if err != nil {
		fmt.Println(err.Error())
	} else {
		response, _ := json.Marshal(*res)
		fmt.Println("Response: " + string(response))
	}
}
