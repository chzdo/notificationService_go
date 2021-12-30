package models

import (
	"time"
)

type UserMobileSettingssModel struct {
	Id        int       `json:"id" bson:"id" `
	CreatedOn time.Time `json:"createdOn"  bson:"createdOn" `
	UpdatedOn time.Time `json:"updatedOn,omitempty"  bson:"updatedOn"`
	IsActive  bool      `json:"isActive" bson:"isActive"`
	IsDeleted bool      `json:"isDeleted" bson:"isDeleted"`
	//main models
	DeviceId  string `json:"deviceId" bson:"deviceId" validate:"required"`
	UserEmail string `json:"userEmail" bson:"userEmail" validate:"required,email"`
}
