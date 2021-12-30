package pushnotification

import (
	"os"

	"github.com/tbalthazar/onesignal-go"
)

type PushNotification struct {
	Driver onesignal.Client
}

func (p *PushNotification) Push(data PushData) error {

	appID := os.Getenv("ONE_SIGNAL_APP_ID")

	notificationReq := &onesignal.NotificationRequest{
		AppID:    appID,
		Contents: map[string]string{"en": data.Message},
		Headings: map[string]string{"en": data.Title},

		IncludePlayerIDs: data.Players,
	}

	_, _, err := p.Driver.Notifications.Create(notificationReq)
	if err != nil {

		return err
	}

	return nil
}

type PushData struct {
	Message string
	Title   string
	Players []string
}
