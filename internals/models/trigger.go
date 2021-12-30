package models

import (
	"time"
)

type model struct{}

type TriggersModel struct {
	Id           int       `json:"id" bson:"id" `
	CreatedOn    time.Time `json:"createdOn"  bson:"createdOn" `
	UpdatedOn    time.Time `json:"updatedOn,omitempty"  bson:"updatedOn"`
	IsActive     bool      `json:"isActive" bson:"isActive"`
	IsDeleted    bool      `json:"isDeleted" bson:"isDeleted"`
	Name         string    `json:"name" bson:"name" validate:"required"`
	Subject      string    `json:"subject" bson:"subject" validate:"required"`
	Type         string    `json:"type" bson:"type" `
	Template     Template  `json:"template" bson:"template" validate:"required"`
	Placeholders []string  `json:"placeholders,-" bson:"-" `
}

type Template struct {
	InAppTemplate string `json:"inAppTemplate" bson:"inAppTemplate" validate:"required"`
	EmailTemplate string `json:"emailTemplate" bson:"emailTemplate"   validate:"required"`
}

type TriggerList []map[string]interface{}

var Triggers = &TriggerList{
	{
		"name": "DEAL_CREATED",
		"type": "organization",
		"placeholders": []string{
			"id", "name", "dealCreated",
		},
	},
	{
		"name": "SIGN_UP",
		"type": "system",
		"placeholders": []string{
			"name", "id", "days",
		},
	},
}

func (t *TriggerList) Get(name string) (int, map[string]interface{}) {

	for key, value := range *t {
		if name == value["name"] {
			return key, value
		}
	}

	return -1, nil
}
