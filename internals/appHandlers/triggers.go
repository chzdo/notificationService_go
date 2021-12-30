package appHandlers

// func (handle *Handler) TriggerRoutes() {

// 	handle.Handler.Get("/triggers", func(rw http.ResponseWriter, r *http.Request) {

// 		response, err := handle.Services.GetTriggersWithoutTemplates(r)

// 		if err != nil {
// 			handle.Responses.ErrorRespond(rw, err)
// 			return
// 		}

// 		handle.Responses.Respond(rw, response)
// 	})

// 	handle.Handler.Get("/triggers/with-template", func(rw http.ResponseWriter, r *http.Request) {

// 		response, err := handle.Services.GetTriggersWithTemplates(r)

// 		if err != nil {
// 			handle.Responses.ErrorRespond(rw, err)
// 			return
// 		}

// 		handle.Responses.Respond(rw, response)
// 	})

// 	handle.Handler.Get("/triggers/with-template/{id}", func(rw http.ResponseWriter, r *http.Request) {

// 		response, err := handle.Services.GetTriggerWithTemplates(r)

// 		if err != nil {
// 			handle.Responses.ErrorRespond(rw, err)
// 			return
// 		}

// 		handle.Responses.Respond(rw, response)
// 	})

// 	handle.Handler.Put("/triggers/with-template/{id}", func(rw http.ResponseWriter, r *http.Request) {

// 		response, err := handle.Services.UpdateTriggersWithTemplates(r)

// 		if err != nil {
// 			handle.Responses.ErrorRespond(rw, err)
// 			return
// 		}

// 		handle.Responses.Respond(rw, response)
// 	})
// 	handle.Handler.Post("/triggers", func(rw http.ResponseWriter, r *http.Request) {

// 		response, err := handle.Services.CreateTriggersWithTemplates(r)

// 		if err != nil {
// 			handle.Responses.ErrorRespond(rw, err)
// 			return
// 		}

// 		handle.Responses.Respond(rw, response)
// 	})

// 	handle.Handler.Delete("/triggers/with-template/{id}", func(rw http.ResponseWriter, r *http.Request) {

// 		response, err := handle.Services.DeleteTriggersWithTemplates(r)

// 		if err != nil {
// 			handle.Responses.ErrorRespond(rw, err)
// 			return
// 		}

// 		handle.Responses.Respond(rw, response)
// 	})
// }
