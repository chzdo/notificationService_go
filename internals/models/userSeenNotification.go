package models

import (
	"time"
)

type UserSeenNotificationModel struct {
	Id        int       `json:"id" bson:"id" `
	CreatedOn time.Time `json:"createdOn"  bson:"createdOn" `
	UpdatedOn time.Time `json:"updatedOn,omitempty"  bson:"updatedOn"`
	IsActive  bool      `json:"isActive" bson:"isActive"`
	IsDeleted bool      `json:"isDeleted" bson:"isDeleted"`

	OrgId           uint   `json:"orgId" bson:"orgId" validate:"required,numeric"`
	UserId          uint   `json:"userId" bson:"userId" validate:"required,numeric"`
	NotificationIds []uint `json:"notificationIds" bson:"notificationIds" validate:"required,dive,required,numeric"`
}
