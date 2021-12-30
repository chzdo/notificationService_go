package models

import "time"

type OrganizationSettingsModel struct {
	Id        int       `json:"id" bson:"id" `
	CreatedOn time.Time `json:"createdOn"  bson:"createdOn" `
	UpdatedOn time.Time `json:"updatedOn,omitempty"  bson:"updatedOn"`
	IsActive  bool      `json:"isActive" bson:"isActive"`
	IsDeleted bool      `json:"isDeleted" bson:"isDeleted"`

	RoleId      int              `json:"roleId" bson:"roleId" validate:"required,numeric"`
	OrgId       int              `json:"orgId" bson:"orgId" validate:"required,numeric"`
	MailingList []MailingList    `json:"mailingList" bson:"mailingList" validate:"required,dive,required"`
	TriggerList []OrgTriggerList `json:"triggerList" bson:"triggerList" validate:"required,dive,required"`
	Recipients  []string         `json:"recipients" bson:"recipients" validate:"required,max=3,dive,required,containsValidRecipients"`
}

type MailingList struct {
	UserId          int    `json:"userId" bson:"userId"  validate:"required,numeric"`
	TeamId          int    `json:"teamId" bson:"teamId"  validate:"required,numeric"`
	SupervisorId    int    `json:"supervisorId" bson:"supervisorId"  validate:"required,numeric"`
	UserEmail       string `json:"userEmail" bson:"userEmail"  validate:"required,email"`
	SupervisorEmail string `json:"supervisorEmail" bson:"supervisorEmail"  validate:"required,email"`
}

type OrgTriggerList struct {
	Name  string `json:"name" bson:"name"  validate:"required,containsValidTriggers"`
	InApp bool   `json:"inApp" bson:"inApp"  validate:"isBoolean"`
	Email bool   `json:"email" bson:"email"  validate:"isBoolean"`
}

var Recipients = map[string]string{
	"all":        "all",
	"initiator":  "initiator",
	"supervisor": "supervisor",
	"team":       "team",
}

func (o OrganizationSettingsModel) CheckRecipients(t []string, key string) bool {

	for _, v := range t {
		if v == key {
			return true
		}
	}
	return false
}
