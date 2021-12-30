package main

import (
	"fmt"
	"net/http"
	"os"
)

func (app *application) serve() error {

	routes, socket := app.Routes()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: routes,
	}

	go socket.Serve()

	defer socket.Close()

	return server.ListenAndServe()

}
