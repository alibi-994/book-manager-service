package main

import (
	"book-manager-server/authenticate"
	"book-manager-server/config"
	"book-manager-server/db"
	"book-manager-server/handlers"
	"log"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	// Read configurations
	var cfg config.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Could not read the DB config")
	}

	// Connect to DB
	gormDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	log.Print("Connected to DB successfully")

	// Migrate Tables
	if err = gormDB.CreateNewSchema(); err != nil {
		log.Fatal("Could not migrate tables")
	}
	log.Print("Tables migrated successfully")

	// Create an instance of authenticate
	auth, err := authenticate.InitAuth(gormDB)

	// Create an instance of bookManagerServer
	bookManagerServer := handlers.BookManagerServer{
		DB:   gormDB,
		Auth: auth,
	}

	// Handle API calls
	http.HandleFunc("/api/v1/auth/signup", bookManagerServer.HandleSignUp)
	http.HandleFunc("/api/v1/auth/login", bookManagerServer.HandleLogin)
	http.HandleFunc("/api/v1/books", bookManagerServer.HandleCrud)
	http.HandleFunc("/api/v1/books/", bookManagerServer.HandleGetBook)

	// Set up http server
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
