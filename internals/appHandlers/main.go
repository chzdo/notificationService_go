package appHandlers

import (
	"notification_service/internals/responses"
	"notification_service/internals/services"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Handler   chi.Router
	Services  *services.Services
	Responses *responses.ResponseFunctions
}
