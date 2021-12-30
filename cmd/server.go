package main

import (
	"fmt"
	"net/http"
)

func (app *application) serve() error {

	routes, socket := app.Routes()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: &routes,
	}

	go socket.Serve()

	defer socket.Close()

	return server.ListenAndServe()

}
