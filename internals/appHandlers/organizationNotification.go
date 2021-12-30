package appHandlers

import (
	"net/http"
)

func (handle *Handler) OrganizationNotificationRoutes() {

	handle.Handler.Handle("/organization-notification/system-single", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		response, err := handle.Services.SendSystemNotificationSingle(r)

		if err != nil {
			handle.Responses.ErrorRespond(rw, err)
			return
		}

		handle.Responses.Respond(rw, response)
	}))

	handle.Handler.Handle("/organization-notification/system-multi", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		response, err := handle.Services.SendSystemNotificationMulti(r)

		if err != nil {
			handle.Responses.ErrorRespond(rw, err)
			return
		}

		handle.Responses.Respond(rw, response)
	})).Methods("POST")

	// handle.Handler.Post("/organization-notification/organization-triggers", func(rw http.ResponseWriter, r *http.Request) {

	// 	response, err := handle.Services.SendOrganizationNotification(r)

	// 	if err != nil {
	// 		handle.Responses.ErrorRespond(rw, err)
	// 		return
	// 	}

	// 	handle.Responses.Respond(rw, response)
	// })

	// handle.Handler.Post("/organization-notification/organization-socials", func(rw http.ResponseWriter, r *http.Request) {

	// 	response, err := handle.Services.SendOrganizationNotificationSocial(r)

	// 	if err != nil {
	// 		handle.Responses.ErrorRespond(rw, err)
	// 		return
	// 	}

	// 	handle.Responses.Respond(rw, response)
	// })
}
