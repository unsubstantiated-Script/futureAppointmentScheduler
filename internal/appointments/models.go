package appointments

import "time"

type Appointment struct {
	ID        int       `json:"id"`
	TrainerID int       `json:"trainer_id"`
	UserID    int       `json:"user_id"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
}

type AvailabilitySlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
