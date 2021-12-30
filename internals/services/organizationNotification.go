package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notification_service/internals/helpers"
	"notification_service/internals/mailing"
	"notification_service/internals/models"
	pushnotification "notification_service/internals/pushNotification"
	"notification_service/internals/responses"
	"notification_service/internals/socket"
	"strconv"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}

func (service *Services) SendOrganizationNotification(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	triggerNotifcation := models.TriggerNotificationModel{}

	err := json.NewDecoder(request.Body).Decode(&triggerNotifcation)

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	varEll := service.Validator.Validate(triggerNotifcation)

	if varEll != nil {
		return nil, responses.SetError(http.StatusBadRequest, varEll[0])
	}

	go service.processTriggerNotification(triggerNotifcation)

	return responses.SetResponse(http.StatusAccepted, "Notification processing"), nil
}

func (service *Services) SendOrganizationNotificationSocial(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	socialTrigger := models.SocialNotificationModel{}

	json.NewDecoder(request.Body).Decode(&socialTrigger)

	errV := service.Validator.Validate(socialTrigger)

	if errV != nil {
		return nil, responses.SetError(http.StatusBadRequest, errV[0])
	}

	metadata := socket.SocketData{
		NameSpace: "/social",
		Room:      fmt.Sprintf("room-%d", socialTrigger.OrgId),
		Event:     fmt.Sprintf("NEW:%s", socialTrigger.Trigger),
		Data:      socialTrigger.Data,
	}

	go service.sendBroadCast(metadata)
	return responses.SetResponse(http.StatusAccepted, "Notification processing"), nil
}
func (service *Services) SendSystemNotificationMulti(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	multiNotification := []models.SingleSystemNotificationModel{}

	json.NewDecoder(request.Body).Decode(&multiNotification)

	valErr := service.Validator.Validate(models.MultiystemNotificationModel{
		Items: multiNotification,
	})

	if valErr != nil {
		return nil, responses.SetError(http.StatusBadRequest, valErr[0])
	}

	service.processMultiNotification(multiNotification, 0, make([]error, 0))

	return responses.SetResponse(http.StatusAccepted, "Notification Processing"), nil
}

func (service *Services) SendSystemNotificationSingle(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	singleNotification := &models.SingleSystemNotificationModel{}

	err := json.NewDecoder(request.Body).Decode(&singleNotification)

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	valErr := service.Validator.Validate(singleNotification)

	if valErr != nil {
		return nil, responses.SetError(http.StatusBadRequest, valErr[0])
	}

	metadata, err := service.processSystemNotification(singleNotification)

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	go service.sendMail(*metadata)

	return responses.SetResponse(http.StatusAccepted, "Notification Processing"), nil

}

func (service *Services) processSystemNotification(s *models.SingleSystemNotificationModel) (*mailing.MailMetaData, error) {

	trigger := models.TriggersModel{}

	query := map[string]interface{}{
		"name":     s.Trigger,
		"isActive": true,
	}

	err := service.getTrigger(query, &trigger)

	if err != nil {
		return nil, err
	}

	if trigger.Type != "system" {

		return nil, fmt.Errorf("Trigger not System Notification")
	}

	metadata := mailing.MailMetaData{
		Type:         trigger.Type,
		Subject:      trigger.Subject,
		Event:        trigger.Name,
		SenderPrefix: "NOTCH-NOTIFICATION",
		Recipients:   []string{s.Email},
		Data:         s.Data,
	}

	return &metadata, nil

}

func (service *Services) processMultiNotification(s []models.SingleSystemNotificationModel, index uint, err []error) []error {

	if len(s) == int(index) {

		return err

	}

	metadata, errs := service.processSystemNotification(&s[index])

	if errs != nil {

		service.Logs.ErrorLogs.Panicln(errs)

		err = append(err, errs)
	}

	go service.sendMail(*metadata)

	index += 1

	return service.processMultiNotification(s, index, err)

}

func (service *Services) sendMail(metadata mailing.MailMetaData) {

	err := service.Mailer.SendMail(metadata)

	if err != nil {

		service.Logs.ErrorLogs.Println(err)
	}

}

func (service *Services) sendPush(metadata pushnotification.PushData) {

	err := service.Push.Push(metadata)

	if err != nil {

		service.Logs.ErrorLogs.Println(err)
	}

}

func (service *Services) sendBroadCast(metadata socket.SocketData) {

	result := service.Socket.BroadCastNotification(metadata.NameSpace, metadata.Room, metadata.Event, metadata.Data)

	service.Logs.ErrorLogs.Println("broadcast state :", result)

}

func (service *Services) processTriggerNotification(t models.TriggerNotificationModel) {

	trigger := models.TriggersModel{}

	query := map[string]interface{}{
		"name":     t.Trigger,
		"isActive": true,
	}
	err := service.getTrigger(query, &trigger)

	if err != nil {
		service.Logs.ErrorLogs.Println(err)
		return
	}

	if trigger.Type != "organization" {
		service.Logs.ErrorLogs.Println("Trigger not organization trigger")
		return
	}

	results, err := service.getDataToSave(t, trigger)
	if err != nil {
		service.Logs.ErrorLogs.Println(err)
		return
	}

	metaData := models.NotificationMetaData{
		MailData: mailing.MailMetaData{
			Subject:      trigger.Subject,
			Event:        trigger.Name,
			Type:         trigger.Type,
			Recipients:   results["emails"].([]string),
			Data:         results["template"].(map[string]interface{}),
			SenderPrefix: fmt.Sprintf("NOTCH FOR %s", t.OrgName),
		},

		PushData: pushnotification.PushData{
			Message: results["template"].(map[string]interface{})["inAppTemplate"].(string),
			Title:   strings.Join(strings.Split(results["trigger"].(string), "_"), " "),
			Players: results["players"].([]string),
		},
		SocketData: socket.SocketData{
			Event:     "NEW_NOTIFICATION",
			Data:      results,
			NameSpace: "/",
			Room:      fmt.Sprintf("room-%d", results["orgId"]),
		},
	}

	go service.sendNotify(metaData)
}

func getSenderSettings(orgSettings []models.OrganizationSettingsModel, creatorId uint) (models.OrganizationSettingsModel, models.MailingList) {
	var senderRoleSettings models.OrganizationSettingsModel
	var sendMailingData models.MailingList
	for _, value := range orgSettings {
		for _, v := range value.MailingList {
			if v.UserId == int(creatorId) {
				senderRoleSettings = value
				sendMailingData = v
			}
		}
	}

	return senderRoleSettings, sendMailingData
}

func setAll(s []models.OrganizationSettingsModel, r *models.Recievers, trigger string) {
	//defer wg.Done()

	localwg := sync.WaitGroup{}
	mailingList := []models.MailingList{}
	triggerList := models.OrgTriggerList{}
	for _, value := range s {

		mailingList = value.MailingList
		localwg.Add(1)
		go func() {

			for _, v := range value.TriggerList {
				if v.Name == trigger {
					triggerList = v
				}
			}

			localwg.Done()
		}()

		localwg.Wait()

		localwg.Add(1)
		go func() {

			for _, k := range mailingList {

				if triggerList.InApp {
					r.InApp[strconv.Itoa(k.UserId)] = uint(k.UserId)
					r.Push[k.UserEmail] = k.UserEmail

				}
				r.Email[k.UserEmail] = k.UserEmail
				r.RoleIds[strconv.Itoa(value.RoleId)] = uint(value.RoleId)

			}
			localwg.Done()
		}()
		localwg.Wait()

	}

	wg.Done()

}

func (service *Services) sendNotify(d models.NotificationMetaData) {

	go service.sendMail(d.MailData)

	go service.sendBroadCast(d.SocketData)

	go service.sendPush(d.PushData)
}

func getTemplate(t *models.OrganizationNotificationModel, placeholders []string) {
	defer wg.Done()

	for _, value := range placeholders {

		t.Template.EmailTemplate = strings.ReplaceAll(t.Template.EmailTemplate, fmt.Sprintf("@%s", value), t.Data[value].(string))
		t.Template.InAppTemplate = strings.ReplaceAll(t.Template.EmailTemplate, fmt.Sprintf("@%s", value), t.Data[value].(string))

	}

}

func setOthers(o models.OrganizationSettingsModel, owner models.MailingList, recievers *models.Recievers, trigger string) {
	defer wg.Done()
	localWg := sync.WaitGroup{}
	triggerList := models.OrgTriggerList{}
	localWg.Add(1)

	go func() {
		for _, value := range o.TriggerList {
			if value.Name == trigger {
				triggerList = value
			}
		}
		localWg.Done()
	}()

	localWg.Wait()

	if o.CheckRecipients(o.Recipients, models.Recipients["team"]) {
		localWg.Add(1)
		go func() {
			for _, value := range o.MailingList {

				if value.TeamId != owner.TeamId {
					continue
				}

				if triggerList.InApp {
					recievers.InApp[strconv.Itoa(value.UserId)] = uint(value.UserId)
					recievers.Push[value.UserEmail] = value.UserEmail

				}
				recievers.Email[value.UserEmail] = value.UserEmail
				recievers.RoleIds[strconv.Itoa(o.RoleId)] = uint(o.RoleId)
			}

			localWg.Done()
		}()
	}

	if o.CheckRecipients(o.Recipients, models.Recipients["initiator"]) {
		localWg.Add(1)
		go func() {

			if triggerList.InApp {
				recievers.InApp[strconv.Itoa(owner.UserId)] = uint(owner.UserId)
				recievers.Push[owner.UserEmail] = owner.UserEmail

			}
			recievers.Email[owner.UserEmail] = owner.UserEmail
			recievers.RoleIds[strconv.Itoa(o.RoleId)] = uint(o.RoleId)

			localWg.Done()
		}()
	}

	if o.CheckRecipients(o.Recipients, models.Recipients["supervisor"]) {
		localWg.Add(1)
		go func() {

			if triggerList.InApp {
				recievers.InApp[strconv.Itoa(owner.SupervisorId)] = uint(owner.SupervisorId)
				recievers.Push[owner.SupervisorEmail] = owner.SupervisorEmail

			}
			recievers.Email[owner.SupervisorEmail] = owner.SupervisorEmail
			recievers.RoleIds[strconv.Itoa(o.RoleId)] = uint(o.RoleId)

			localWg.Done()
		}()
	}

	localWg.Wait()

}

func (service *Services) getRecievers(t models.TriggerNotificationModel, trigger models.TriggersModel) models.Recievers {
	query := map[string]interface{}{
		"orgId":            t.OrgId,
		"triggerList.name": trigger.Name,
	}

	orgSettings := []models.OrganizationSettingsModel{}
	service.getOrganizationSettings(query, &orgSettings)

	recievers := models.Recievers{
		Email:   make(map[string]string),
		InApp:   make(map[string]uint),
		Push:    make(map[string]string),
		RoleIds: make(map[string]uint),
	}

	senderRoleSetting, ownerList := getSenderSettings(orgSettings, t.CreatorId)

	if check := senderRoleSetting.CheckRecipients(senderRoleSetting.Recipients, "all"); check == true {
		wg.Add(1)
		go setAll(orgSettings, &recievers, trigger.Name)
		wg.Wait()

	} else {
		wg.Add(1)
		go setOthers(senderRoleSetting, ownerList, &recievers, trigger.Name)
		wg.Wait()

	}

	return recievers
}

func (service *Services) getDataToSave(t models.TriggerNotificationModel, trigger models.TriggersModel) (map[string]interface{}, error) {

	recievers := service.getRecievers(t, trigger)

	dataToSave := models.OrganizationNotificationModel{
		OrgId:     t.OrgId,
		CreatorId: t.CreatorId,
		Trigger:   trigger.Name,
		Data:      t.Data,
		Template:  trigger.Template,
		RoleIds:   helpers.GetMapUintKeys(recievers.RoleIds),
		Recievers: helpers.GetMapUintKeys(recievers.InApp),
	}

	emails := helpers.GetMapStringstKeys(recievers.Email)
	pushMails := helpers.GetMapStringstKeys(recievers.Push)

	players, err := service.getPlayers(pushMails)

	if err != nil {

		return nil, err
	}

	wg.Add(1)
	go getTemplate(&dataToSave, trigger.Placeholders)
	wg.Wait()

	results, err := service.Models["organizationnotifications"].Insert(dataToSave)

	if err != nil {
		return nil, err
	}
	results["emails"] = emails
	results["players"] = players
	return results, nil
}

func (service *Services) getOrganizationNotification(query map[string]interface{}) ([]models.OrganizationNotificationModel, error) {

	results, err := service.Models["organizationnotifications"].GetAll(query)

	if err != nil {
		return nil, err
	}
	itemToReturn := []models.OrganizationNotificationModel{}
	err = json.Unmarshal(results, &itemToReturn)
	if err != nil {
		return nil, err
	}
	return itemToReturn, nil
}
