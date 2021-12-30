package services

import (
	"notification_service/internals/appValidators"
	"notification_service/internals/logger"
	"notification_service/internals/mailing"
	"notification_service/internals/models"
	pushnotification "notification_service/internals/pushNotification"
	"notification_service/internals/socket"
)

type Services struct {
	Logs      logger.Logger
	Models    map[string]*models.DBModels
	Validator appValidators.ValidatorStruct
	Socket    *socket.Socket
	Mailer    mailing.Mailer
	Push      pushnotification.PushNotification
}
