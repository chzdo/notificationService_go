package appHandlers

import "net/http"

func (handle *Handler) UserSeenNotificationRoutes() {

	handle.Handler.Post("/user-seen-notification", func(rw http.ResponseWriter, r *http.Request) {

		response, err := handle.Services.CreateUserSeenNotification(r)

		if err != nil {
			handle.Responses.ErrorRespond(rw, err)
			return
		}

		handle.Responses.Respond(rw, response)
	})

	handle.Handler.Get("/user-seen-notification/{orgId}/{userId}", func(rw http.ResponseWriter, r *http.Request) {

		response, err := handle.Services.GetUserSeenNotification(r)

		if err != nil {
			handle.Responses.ErrorRespond(rw, err)
			return
		}

		handle.Responses.Respond(rw, response)
	})

}
