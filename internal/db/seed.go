package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

// SeedAppointment maps the delivered seed file keys. The take-home instructions
// require keeping `started_at` and `ended_at` in JSON, even though the DB/domain
// columns use `starts_at` and `ends_at`.
type SeedAppointment struct {
	ID        int    `json:"id"`
	TrainerID int    `json:"trainer_id"`
	UserID    int    `json:"user_id"`
	StartsAt  string `json:"started_at"`
	EndsAt    string `json:"ended_at"`
}

// RunMigrations runs the migrations
func RunMigrations(db *sql.DB) error {
	sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read sql file: %w", err)
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("failed to run sql: %w", err)
	}
	return nil
}

// SeedAppointments inserts the appointments from the json file
func SeedAppointments(db *sql.DB) error {
	//Making sure we've not already seeded
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM appointments").Scan(&count); err != nil {
		return fmt.Errorf("failed to count appointments: %w", err)
	}

	if count > 0 {
		return nil
	}

	//Reading the json file
	bytes, err := os.ReadFile("data/appointments.json")
	if err != nil {
		return fmt.Errorf("failed to read appointments json: %w", err)
	}

	//Unmarshalling the json
	var appointments []SeedAppointment
	if err := json.Unmarshal(bytes, &appointments); err != nil {
		return fmt.Errorf("failed to unmarshal appointments json: %w", err)
	}

	//Inserting the appointments
	ctx := context.Background()

	for _, appointment := range appointments {
		_, err := db.ExecContext(ctx, "INSERT INTO appointments (id, trainer_id, user_id, starts_at, ends_at) VALUES ($1, $2, $3, $4, $5)",
			appointment.ID,
			appointment.TrainerID,
			appointment.UserID,
			appointment.StartsAt,
			appointment.EndsAt)
		if err != nil {
			return fmt.Errorf("failed to insert appointment: %w", err)
		}
	}

	//Setting the sequence
	if _, err := db.Exec(`
		SELECT setval(
			   pg_get_serial_sequence('appointments', 'id'),
			   COALESCE((SELECT MAX(id) FROM appointments), 1),
			   true
		)
`); err != nil {
		return fmt.Errorf("failed to set sequence: %w", err)
	}

	return nil
}
