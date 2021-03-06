package mailing

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type Mailer struct {
	Driver mailgun.Mailgun
}

func (mailer *Mailer) SendMail(metadata MailMetaData) error {

	sender := os.Getenv("APP_DEFAULT_EMAIL")

	from := metadata.SenderPrefix + fmt.Sprintf(" <%s>", sender)
	mg := mailer.Driver.NewMessage(from, metadata.Subject, "", metadata.Recipients...)

	var templateString string = metadata.Event

	if metadata.Type == "organization" {
		templateString = "GENERAL"
	}

	tmpl, err := template.ParseFiles(fmt.Sprintf("internals/mailing/templates/%s.tmpl", templateString))

	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	tmpl.Execute(body, metadata.Data)

	mg.SetHtml(body.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, _, err = mailer.Driver.Send(ctx, mg)

	return err

}

type MailMetaData struct {
	Subject      string
	SenderPrefix string
	Event        string
	Type         string
	Recipients   []string
	Data         map[string]interface{}
}
