package appHandlers

import (
	"notification_service/internals/responses"
	"notification_service/internals/services"

	"github.com/gorilla/mux"
)

type Handler struct {
	Handler   mux.Router
	Services  *services.Services
	Responses *responses.ResponseFunctions
}
