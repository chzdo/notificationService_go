package main

import (
	"net/http"
	"notification_service/internals/logger"
)

type config struct {
	port int
	env  string
	db   struct {
		DB_URI string
	}
}

type application struct {
	config  *config
	handler http.Handler
	logs    *logger.Logger
}
