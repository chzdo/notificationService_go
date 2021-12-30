package appValidators

import (
	"encoding/json"
	"fmt"

	"notification_service/internals/models"

	"github.com/go-playground/validator/v10"
)

type ValidatorStruct struct {
	validator validator.Validate
}

var Validator = &ValidatorStruct{
	validator: *validator.New(),
}

func (t *ValidatorStruct) Validate(s interface{}) []string {

	t.validator.RegisterValidation("containsValidRecipients", func(fl validator.FieldLevel) bool {
		t := fl.Field().String()

		if models.Recipients[t] == "" {
			return false
		}

		return true
	})

	t.validator.RegisterValidation("containsValidTriggers", func(fl validator.FieldLevel) bool {
		t := fl.Field().String()

		if index, _ := models.Triggers.Get(t); index == -1 {
			return false
		}

		return true
	})

	t.validator.RegisterValidation("isBoolean", func(fl validator.FieldLevel) bool {
		_, ok := fl.Field().Interface().(bool)

		if !ok {
			return false
		}

		return true
	})

	t.validator.RegisterValidation("isValidPlaceholders", func(fl validator.FieldLevel) bool {

		v := fl.Parent().Interface()

		byte, err := json.Marshal(v)

		if err != nil {
			return false
		}

		temp := map[string]interface{}{}

		json.Unmarshal(byte, &temp)

		index, value := models.Triggers.Get(temp["trigger"].(string))

		if index < 0 {
			return false
		}
		dataTemp := temp["data"].(map[string]interface{})

		for _, value := range value["placeholders"].([]string) {

			if dataTemp[value] == nil {
				return false
			}
		}

		return true
	})

	err := t.validator.Struct(s)

	var sliceToReturn []string

	if err == nil {

		return sliceToReturn
	}

	for _, err := range err.(validator.ValidationErrors) {

		sliceToReturn = append(sliceToReturn, fmt.Sprintf("%s failed %s validation", err.Field(), err.Tag()))

	}
	return sliceToReturn
}
