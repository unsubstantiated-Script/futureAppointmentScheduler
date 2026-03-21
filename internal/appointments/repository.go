package appointments

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByTrainer(trainerID int) ([]Appointment, error) {
	//Querying the database
	rows, err := r.db.Query(`
		SELECT id, trainer_id, user_id, starts_at, ends_at
		FROM appointments
		WHERE trainer_id = $1
		ORDER BY starts_at
	`, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Building the appointments
	var appointments []Appointment

	for rows.Next() {
		var a Appointment
		if err := rows.Scan(
			&a.ID,
			&a.TrainerID,
			&a.UserID,
			&a.StartsAt,
			&a.EndsAt,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *Repository) CreateAppointment(appt Appointment) error {
	_, err := r.db.Exec(`
		INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at)
		VALUES ($1, $2, $3, $4)`,
		appt.TrainerID,
		appt.UserID,
		appt.StartsAt,
		appt.EndsAt,
	)
	return err
}

func (r *Repository) GetBetween(trainerID int, startsAt, endsAt time.Time) ([]Appointment, error) {
	//Querying the database
	rows, err := r.db.Query(`
		SELECT id, trainer_id, user_id, starts_at, ends_at
		FROM appointments
		WHERE trainer_id = $1 
			AND starts_at < $3 
			AND ends_at > $2 
		ORDER BY starts_at
`, trainerID, startsAt, endsAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Building the appointments
	var appointments []Appointment

	for rows.Next() {
		var a Appointment
		if err := rows.Scan(
			&a.ID,
			&a.TrainerID,
			&a.UserID,
			&a.StartsAt,
			&a.EndsAt,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}
