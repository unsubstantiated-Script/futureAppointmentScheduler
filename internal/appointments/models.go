package appointments

import "time"

type Appointment struct {
	ID        int       `json:"id"`
	TrainerID int       `json:"trainer_id"`
	UserID    int       `json:"user_id"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
}
