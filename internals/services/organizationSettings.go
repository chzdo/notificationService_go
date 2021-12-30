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

func (service *Services) CreateOrganizationSettings(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	organizationSettings := &models.OrganizationSettingsModel{}

	err := json.NewDecoder(request.Body).Decode(&organizationSettings)

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	errVal := service.Validator.Validate(organizationSettings)

	if errVal != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, errVal[0])
	}

	query := map[string]interface{}{
		"orgId":  organizationSettings.OrgId,
		"roleId": organizationSettings.RoleId,
	}
	_, err = service.Models["organizationsettings"].Get(query)

	if err == nil {
		return nil, responses.SetError(http.StatusNotAcceptable, "settings for this org and role has been created already")
	}

	result, err := service.Models["organizationsettings"].Insert(organizationSettings)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	return responses.SetResponse(http.StatusCreated, result), nil
}

func (service *Services) GetOrganizationSettings(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	organizationSettings := &models.OrganizationSettingsModel{}

	orgId, err := strconv.Atoi(chi.URLParam(request, "orgId"))

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	roleId, err := strconv.Atoi(chi.URLParam(request, "roleId"))

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	query := map[string]interface{}{
		"orgId":  orgId,
		"roleId": roleId,
	}

	result, err := service.Models["organizationsettings"].Get(query)

	if err != nil {

		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	err = json.Unmarshal(result, &organizationSettings)

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	return responses.SetResponse(http.StatusOK, organizationSettings), nil

}
func (service *Services) UpdateOrganizationSettings(request *http.Request) (*responses.SuccessResponse, *responses.ErrorResponse) {

	organizationSettings := models.OrganizationSettingsModel{}

	orgId, err := strconv.Atoi(chi.URLParam(request, "orgId"))

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	roleId, err := strconv.Atoi(chi.URLParam(request, "roleId"))

	if err != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, err.Error())
	}

	bytebody, err := io.ReadAll(request.Body)

	if err != nil {
		return nil, responses.SetError(http.StatusBadRequest, err.Error())
	}

	err = json.Unmarshal(bytebody, &organizationSettings)

	organizationSettings.OrgId = orgId
	organizationSettings.RoleId = roleId

	valErr := service.Validator.Validate(organizationSettings)

	if valErr != nil {
		return nil, responses.SetError(http.StatusUnprocessableEntity, valErr[0])
	}

	query := map[string]interface{}{
		"orgId":  orgId,
		"roleId": roleId,
	}

	updatebody := map[string]interface{}{}

	err = json.Unmarshal(bytebody, &updatebody)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}
	result, err := service.Models["organizationsettings"].Update(query, updatebody)

	if err != nil {
		return nil, responses.SetError(http.StatusInternalServerError, err.Error())
	}

	if result.MatchedCount < 1 {
		return nil, responses.SetError(http.StatusNotFound, "Resource not found")
	}

	if result.ModifiedCount < 1 {
		return nil, responses.SetError(http.StatusNotModified, "No Update")
	}

	return responses.SetResponse(http.StatusOK, updatebody), nil
}

func (service *Services) getOrganizationSettings(query map[string]interface{}, t *[]models.OrganizationSettingsModel) error {

	result, err := service.Models["organizationsettings"].GetAll(query)

	if err != nil {

		return err
	}

	err = json.Unmarshal(result, t)

	if err != nil {
		return err
	}

	return nil
}
