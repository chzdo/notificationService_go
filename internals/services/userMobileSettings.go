package services

import (
	"encoding/json"
	"net/http"
	"notification_service/internals/models"
	"notification_service/internals/responses"
)

func (service *Services) CreateUserMobileSettings(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	userMobileSettings := &models.UserMobileSettingssModel{}

	err := json.NewDecoder(request.Body).Decode(&userMobileSettings)

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	valErr := service.Validator.Validate(userMobileSettings)

	if valErr != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, valErr[0])
	}

	query := map[string]interface{}{

		"deviceId": userMobileSettings.DeviceId,
	}

	_, err = service.Models["usermobilesettings"].Get(query)

	if err == nil {
		return nil, responses.SetError(http.StatusBadRequest, "This device is already registered")
	}

	result, err := service.Models["usermobilesettings"].Insert(userMobileSettings)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	// service.Socket.BroadCastNotification("/", "orgId", "NEW:TWEET", result)

	// 	err = service.Mailer.SendMail([]string{"chido.nduaguibe@gmail.com"}, mailing.MailMetaData{
	// 	Event:        "DEAL",
	// 	Type:         "organization",
	// 	Subject:      "test",
	// 	SenderPrefix: "NOTCh",
	// })

	// fmt.Println(err)

	// 	err = service.Push.Push(pushnotification.PushData{
	// 	Message: "test",
	// 	Title:   "test",
	// 	Players: []string{"340de08c-6828-11ec-8b20-6a7d14bdd1f5"},
	// })

	return responses.SetResponse(http.StatusCreated, result), nil
}

func (service *Services) getPlayers(email []string) ([]string, error) {

	temp := make([]models.UserMobileSettingssModel, 0)

	playerTemp := []string{}
	query := map[string]interface{}{
		"userEmail": map[string]interface{}{
			"$in": email,
		},
	}
	bytes, err := service.Models["usermobilesettings"].GetAll(query)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &temp)
	if err != nil {
		return nil, err
	}

	for index, _ := range temp {
		playerTemp = append(playerTemp, temp[index].DeviceId)
	}

	return playerTemp, nil
}
