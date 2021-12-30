package models

import (
	"notification_service/internals/mailing"
	pushnotification "notification_service/internals/pushNotification"
	"notification_service/internals/socket"

	"time"
)

type OrganizationNotificationModel struct {
	Id        int       `json:"id" bson:"id" `
	CreatedOn time.Time `json:"createdOn"  bson:"createdOn" `
	UpdatedOn time.Time `json:"updatedOn,omitempty"  bson:"updatedOn"`
	IsActive  bool      `json:"isActive" bson:"isActive"`
	IsDeleted bool      `json:"isDeleted" bson:"isDeleted"`

	OrgId     uint                   `json:"orgId" bson:"orgId" validate:"required,numeric"`
	RoleIds   []uint                 `json:"roleIds" bson:"roleIds" validate:"required,dive,required, numeric"`
	CreatorId uint                   `json:"creatorId" bson:"creatorId" validate:"required,numeric"`
	Recievers []uint                 `json:"recievers" bson:"recievers"  validate:"required, dive, required,numeric"`
	Template  Template               `json:"template" bson:"template" validate:"required,dive, required"`
	Trigger   string                 `json:"trigger" bson:"trigger"  validate:"required" `
	Data      map[string]interface{} `json:"data" bson:"data" validate:"required"`
}

type SingleSystemNotificationModel struct {
	Trigger string                 `json:"trigger" validate:"required,containsValidTriggers"`
	Email   string                 `json:"email" validate:"required,email"`
	Data    map[string]interface{} `json:"data" validate:"required,isValidPlaceholders"`
}

type TriggerNotificationModel struct {
	Trigger   string                 `json:"trigger" validate:"required,containsValidTriggers"`
	OrgId     uint                   `json:"orgId"   validate:"required,numeric"`
	OrgName   string                 `json:"orgName"   validate:"required"`
	CreatorId uint                   `json:"creatorId"   validate:"required,numeric"`
	Data      map[string]interface{} `json:"data"    validate:"required,isValidPlaceholders"`
}

type SocialNotificationModel struct {
	Trigger string                 `json:"trigger" validate:"required"`
	OrgId   uint                   `json:"orgId" validate:"required,numeric"`
	Data    map[string]interface{} `json:"data" validate:"required,min=1"`
}

type MultiystemNotificationModel struct {
	Items []SingleSystemNotificationModel `validate:"required,dive,required"`
}

type Recievers struct {
	Email   map[string]string
	InApp   map[string]uint
	Push    map[string]string
	RoleIds map[string]uint
}

type NotificationMetaData struct {
	MailData   mailing.MailMetaData
	PushData   pushnotification.PushData
	SocketData socket.SocketData
}
