package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/internal/data"
	"server/internal/driver"
)

// config struct holds application config
type config struct {
	port int
}

// application struct holds all configurations that needs to be shared globally across the api
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
	environment string
}

func main() {
	var cfg config
	cfg.port = 8092

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")

	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models:   data.New(db.SQL),
		environment: environment,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	app.infoLog.Println("Server listening on PORT", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
