package main

import (
	"futureAppointmentScheduler/internal/appointments"
	"futureAppointmentScheduler/internal/db"
	"log"
	"net/http"
	"os"
)

func main() {
	//Setting up the API
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//Setting up the DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	dbConn, err := db.Open(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := db.RunMigrations(dbConn); err != nil {
		log.Fatal(err)
	}

	//Building routes
	mux := http.NewServeMux()
	h := appointments.NewHandler()
	mux.HandleFunc("/appointments", h.Appointments)
	mux.HandleFunc("/availability", h.Availability)

	addr := ":" + port
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
