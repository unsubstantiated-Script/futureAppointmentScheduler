package appointments

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetByTrainer returns all the appointments for a given trainer
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

// CreateAppointment inserts a new appointment into the database and returns its ID.
func (r *Repository) CreateAppointment(appt Appointment) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO appointments (trainer_id, user_id, starts_at, ends_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		appt.TrainerID,
		appt.UserID,
		appt.StartsAt,
		appt.EndsAt,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == "23P01" &&
			pgErr.ConstraintName == "appointments_no_overlap" {
			return 0, ErrAppointmentOverlap
		}
		return 0, err
	}

	return id, nil
}

// GetBetween returns all the appointments for a given trainer between two dates
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
