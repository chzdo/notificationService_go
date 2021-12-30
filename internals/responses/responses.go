package responses

import (
	"encoding/json"
	"net/http"
	"notification_service/internals/logger"
)

type SuccessResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ResponsesFunc interface {
	ErrorRespondFunc(http.ResponseWriter, ErrorResponse)
	Respond(response http.ResponseWriter, data SuccessResponse)
}

func (rfs *ResponseFunctions) Respond(response http.ResponseWriter, data *SuccessResponse) {
	js, _ := json.MarshalIndent(data, "", " ")
	response.WriteHeader(data.Status)
	response.Write(js)
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (rfs *ResponseFunctions) ErrorRespond(response http.ResponseWriter, data *ErrorResponse) {
	js, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		rfs.Logs.ErrorLogs.Println(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	rfs.Logs.ErrorLogs.Println(data.Message)
	response.WriteHeader(data.Status)
	response.Write(js)

}

type ResponseFunctions struct {
	Logs *logger.Logger
}

func SetError(code int, message string) *ErrorResponse {

	if message == "mongo: no documents in result" {
		code = http.StatusNotFound
		message = "Resource Not found"
	}
	return &ErrorResponse{
		Status:  code,
		Success: false,
		Message: message,
	}
}

func SetResponse(code int, data interface{}) *SuccessResponse {

	return &SuccessResponse{
		Status:  code,
		Success: true,
		Data:    data,
	}
}
