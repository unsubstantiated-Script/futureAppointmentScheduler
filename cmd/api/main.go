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

// main initializes and starts the application, setting up the database, running migrations, and configuring HTTP routes and server.
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

	// Establishing Database Connection
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

	if err := db.RunMigrations(dbConn); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ran migrations")

	if err := db.SeedAppointments(dbConn); err != nil {
		log.Fatal(err)
	}
	fmt.Println("seeded appointments")

	//Instantiating repo, service, and handler
	repo := appointments.NewRepository(dbConn)
	svc := appointments.NewService(repo)
	handler := appointments.NewHandler(svc)

	//Building routes
	mux := http.NewServeMux()
	mux.HandleFunc("/appointments", handler.Appointments)
	mux.HandleFunc("/availability", handler.Availability)

	addr := ":" + port
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
