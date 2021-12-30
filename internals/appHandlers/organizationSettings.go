package appHandlers

import (
	"net/http"
)

func (handle *Handler) OrganizationSettingsRoutes() {

	handle.Handler.Post("/organization-settings", func(rw http.ResponseWriter, r *http.Request) {

		response, err := handle.Services.CreateOrganizationSettings(r)

		if err != nil {
			handle.Responses.ErrorRespond(rw, err)
			return
		}

		handle.Responses.Respond(rw, response)
	})

	handle.Handler.Get("/organization-settings/{orgId}/{roleId}", func(rw http.ResponseWriter, r *http.Request) {

		response, err := handle.Services.GetOrganizationSettings(r)

		if err != nil {
			handle.Responses.ErrorRespond(rw, err)
			return
		}

		handle.Responses.Respond(rw, response)
	})

	handle.Handler.Put("/organization-settings/{orgId}/{roleId}", func(rw http.ResponseWriter, r *http.Request) {

		response, err := handle.Services.UpdateOrganizationSettings(r)

		if err != nil {
			handle.Responses.ErrorRespond(rw, err)
			return
		}

		handle.Responses.Respond(rw, response)
	})

}
