package utils

import (
	"NeraJima/configs"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var apiKey, emailSender string = configs.SendGridKeyAndFrom()
var client *sendgrid.Client = sendgrid.NewSendClient(apiKey)

func SendRegistrationEmail(name, email string, code int) {
	from := mail.NewEmail("NeraJima", emailSender)
	tos := []*mail.Email{ // list of emails to send this email to
		mail.NewEmail(name, email),
	}

	m := mail.NewV3Mail()
	m.SetFrom(from)
	m.SetTemplateID("d-bccd2db8db3e4699b3e636b78bddb90e")

	p := mail.NewPersonalization()
	p.SetDynamicTemplateData("full_name", name)
	p.SetDynamicTemplateData("verification_code", code)
	p.AddTos(tos...)

	m.AddPersonalizations(p)

	res, err := client.Send(m)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.StatusCode)
	}
}
