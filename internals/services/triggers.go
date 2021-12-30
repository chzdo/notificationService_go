package services

import (
	"encoding/json"
	"net/http"
	"notification_service/internals/models"
	"notification_service/internals/responses"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (service *Services) CreateTriggersWithTemplates(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	trigger := models.TriggersModel{}

	err := json.NewDecoder(request.Body).Decode(&trigger)

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	valErr := service.Validator.Validate(trigger)

	if valErr != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, valErr[0])
	}

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	index, triggerInfo := models.Triggers.Get(trigger.Name)

	if index < 0 {

		return nil, responses.SetError(http.StatusUnprocessableEntity, "Trigger Not found")
	}

	_, err = service.Models["triggers"].Get(map[string]interface{}{
		"name":     trigger.Name,
		"isActive": true,
	})

	if err == nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, "Trigger Already Exist")
	}

	trigger.Type = triggerInfo["type"].(string)

	result, err := service.Models["triggers"].Insert(trigger)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	return responses.SetResponse(http.StatusCreated, result), nil

}

func (service *Services) GetTriggersWithTemplates(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	triggers := []models.TriggersModel{}

	query := map[string]interface{}{
		"isActive": true,
	}

	result, err := service.Models["triggers"].GetAll(query)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	if err = json.Unmarshal(result, &triggers); err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	for index, value := range triggers {
		_, trigger := models.Triggers.Get(value.Name)
		triggers[index].Placeholders = trigger["placeholders"].([]string)
	}

	return responses.SetResponse(http.StatusOK, triggers), nil
}

func (service *Services) GetTriggerWithTemplates(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	triggers := models.TriggersModel{}

	id, err := strconv.Atoi(chi.URLParam(request, "id"))

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	var query = map[string]interface{}{
		"id":       id,
		"isActive": true,
	}

	err = service.getTrigger(query, &triggers)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	return responses.SetResponse(http.StatusOK, triggers), nil

}

func (service *Services) GetTriggersWithoutTemplates(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	return responses.SetResponse(http.StatusOK, models.Triggers), nil
}

func (service *Services) UpdateTriggersWithTemplates(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	updateBody := map[string]interface{}{}

	err = json.NewDecoder(request.Body).Decode(&updateBody)

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	var query = map[string]interface{}{
		"id":       id,
		"isActive": true,
	}

	result, err := service.Models["triggers"].Update(query, updateBody)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	if result.MatchedCount < 1 {
		return nil, responses.SetError(http.StatusNotFound, "Trigger not found")
	}

	if result.ModifiedCount < 1 {
		return nil, responses.SetError(http.StatusNotModified, "No Update")
	}

	updateBody["id"] = id

	return responses.SetResponse(http.StatusOK, updateBody), nil

}
func (service *Services) DeleteTriggersWithTemplates(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	var query = map[string]interface{}{
		"id":       id,
		"isActive": true,
	}

	result, err := service.Models["triggers"].Delete(query)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	if result.MatchedCount < 1 {
		return nil, responses.SetError(http.StatusNotFound, "Trigger not found")
	}

	if result.ModifiedCount < 1 {
		return nil, responses.SetError(http.StatusNotModified, "Not Deleted")
	}

	return responses.SetResponse(http.StatusOK, "Trigger Deleted"), nil

}

func (service *Services) getTrigger(query map[string]interface{}, triggers *models.TriggersModel) error {
	result, err := service.Models["triggers"].Get(query)

	if err != nil {

		return err
	}

	if err = json.Unmarshal(result, &triggers); err != nil {
		return err
	}

	_, trigger := models.Triggers.Get(triggers.Name)

	triggers.Placeholders = trigger["placeholders"].([]string)

	return nil
}
