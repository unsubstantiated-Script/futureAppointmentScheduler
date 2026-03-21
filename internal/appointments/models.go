package appointments

import "time"

type Appointment struct {
	ID        int
	TrainerID int
	UserID    int
	StartsAt  time.Time
	EndsAt    time.Time
}
