package appHandlers

import (
	"notification_service/internals/responses"
	"notification_service/internals/services"

	"github.com/go-chi/chi"
)

type Handler struct {
	Handler   *chi.Mux
	Services  *services.Services
	Responses *responses.ResponseFunctions
}
