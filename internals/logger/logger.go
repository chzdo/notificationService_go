package logger

import (
	"log"
	"os"
)

type Logger struct {
	ErrorLogs *log.Logger
	InfoLogs  *log.Logger
}

func (logs *Logger) Set(env string) {

	var errorLog *log.Logger
	var infoLog *log.Logger
	if env == "production" {
		errorFile, err := os.Open("errorlog.logs")

		if err != nil {
			log.Panicln(err)
		}

		infoFile, err := os.Open("infolog.logs")

		errorLog = log.New(errorFile, "ERROR\t", log.Ldate|log.LUTC|log.Lshortfile)
		infoLog = log.New(infoFile, "INFO\t", log.Ldate|log.LUTC)

	} else {
		errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Lshortfile)
		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)
	}

	logs.ErrorLogs = errorLog
	logs.InfoLogs = infoLog
}
