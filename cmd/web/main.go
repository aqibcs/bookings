package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/aqibcs/bookings/internal/config"
	"github.com/aqibcs/bookings/internal/driver"
	"github.com/aqibcs/bookings/internal/handlers"
	"github.com/aqibcs/bookings/internal/helpers"
	"github.com/aqibcs/bookings/internal/models"
	"github.com/aqibcs/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatalf("Error during application startup: %s", err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatalf("Error starting server: %s", err)
	}
}

func run() (*driver.DB, error) {
	// Register types for gob encoding
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})

	// Set up application configuration
	app.InProduction = false

	// Set up loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Set up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// Construct the DSN string
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		config.Host, config.Port, config.Username, config.DatabaseName, config.Password)

	// Connect to the database
	infoLog.Println("Connecting to the database...")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		errorLog.Fatalf("Cannot connect to the database: %s", err)
	}

	// Create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		errorLog.Fatalf("Cannot create template cache: %s", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	// Set up repository, handlers, renderer, and helpers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
