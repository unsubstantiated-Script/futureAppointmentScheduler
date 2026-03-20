package main

import (
	"database/sql"
	"fmt"
	"futureAppointmentScheduler/internal/appointments"
	"futureAppointmentScheduler/internal/db"
	"log"
	"net/http"
	"os"
	"time"
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

	var dbConn *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		dbConn, err = db.Open(dbURL)
		if err == nil {
			err = dbConn.Ping()
		}
		if err == nil {
			break
		}

		log.Printf("waiting for db... attempt %d: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbConn.Close()

	log.Println("connected to db")

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to db")

	if err := db.RunMigrations(dbConn); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ran migrations")

	if err := db.SeedAppointments(dbConn); err != nil {
		log.Fatal(err)
	}
	fmt.Println("seeded appointments")

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
