package main

import (
	"flag"
	"log"
	"notification_service/internals/logger"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func init() {

	godotenv.Load()
}

var wg sync.WaitGroup

func main() {

	var cfg config

	flag.StringVar(&cfg.env, "env", "development", "app mode")
	flag.IntVar(&cfg.port, "port", 3700, "app port")

	flag.Parse()

	logs := &logger.Logger{}

	logs.Set(cfg.env)

	app := &application{
		config: &cfg,
		logs:   logs,
	}

	log.Printf("running on %s\n", os.Getenv("PORT"))
	err := app.serve()

	if err != nil {
		log.Panic(err)
	}

}
