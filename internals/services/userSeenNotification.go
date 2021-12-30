package services

import (
	"encoding/json"
	"io"
	"net/http"
	"notification_service/internals/models"
	"notification_service/internals/responses"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (service *Services) CreateUserSeenNotification(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	userSeenNotification := models.UserSeenNotificationModel{}

	bytes, err := io.ReadAll(request.Body)

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	json.Unmarshal(bytes, &userSeenNotification)

	valErr := service.Validator.Validate(userSeenNotification)

	if valErr != nil {
		return nil, responses.SetError(http.StatusBadRequest, valErr[0])
	}

	query := map[string]interface{}{
		"orgId":  userSeenNotification.OrgId,
		"userId": userSeenNotification.UserId,
	}

	updateBody := map[string]interface{}{}
	err = json.Unmarshal(bytes, &updateBody)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	_, err = service.Models["userseennotification"].UpdateOrInsert(query, updateBody)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	return responses.SetResponse(http.StatusOK, updateBody), nil
}

func (service *Services) GetUserSeenNotification(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	orgId, err := strconv.Atoi(chi.URLParam(request, "orgId"))

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	userId, err := strconv.Atoi(chi.URLParam(request, "userId"))

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}
	query := map[string]interface{}{
		"orgId":  orgId,
		"userId": userId,
	}

	results, err := service.Models["userseennotification"].Get(query)

	query = map[string]interface{}{
		"orgId":     orgId,
		"recievers": userId,
	}
	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	notifications, err := service.getOrganizationNotification(query)

	seenNotification := models.UserSeenNotificationModel{}
	json.Unmarshal(results, &seenNotification)

	itemToSend := map[string]interface{}{
		"notifications":    notifications,
		"seenNotification": seenNotification,
	}
	return responses.SetResponse(http.StatusOK, itemToSend), nil
}
