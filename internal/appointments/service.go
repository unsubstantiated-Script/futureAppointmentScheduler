package appointments

import (
	"errors"
	"time"
	_ "time/tzdata"
)

var (
	ErrInvalidAppointment = errors.New("invalid appointment")
	ErrAppointmentOverlap = errors.New("appointment overlaps")
)

type Service struct {
	repo *Repository
}

// NewService returns a new Service instance
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetByTrainer returns all the appointments for a given trainer
func (s *Service) GetByTrainer(trainerID int) ([]Appointment, error) {
	return s.repo.GetByTrainer(trainerID)
}

// CreateAppointment inserts a new appointment into the database and returns its ID.
func (s *Service) CreateAppointment(appt Appointment) (int, error) {
	if err := validateAppointment(appt); err != nil {
		return 0, err
	}

	existing, err := s.repo.GetBetween(appt.TrainerID, appt.StartsAt, appt.EndsAt)
	if err != nil {
		return 0, err
	}

	if len(existing) > 0 {
		return 0, ErrAppointmentOverlap
	}

	return s.repo.CreateAppointment(appt)
}

// GetAvailability returns a list of available slots for a given trainer between two dates
func (s *Service) GetAvailability(trainerID int, startsAt, endsAt time.Time) ([]AvailabilitySlot, error) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}

	if !startsAt.Before(endsAt) {
		return nil, ErrInvalidAppointment
	}

	existing, err := s.repo.GetBetween(trainerID, startsAt, endsAt)
	if err != nil {
		return nil, err
	}

	var slots []AvailabilitySlot

	startPT := startsAt.In(loc)
	endPT := endsAt.In(loc)

	// Iterate over each day in the range
	for day := dayStart(startPT); day.Before(endPT); day = day.AddDate(0, 0, 1) {
		if !isBusinessDay(day) {
			continue
		}

		// Iterate over each hour in the day
		for _, hour := range []int{8, 9, 10, 11, 12, 13, 14, 15, 16} {

			// Iterate over each minute in the hour
			for _, minute := range []int{0, 30} {
				slotStartPT := time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, loc)
				slotEndPT := slotStartPT.Add(30 * time.Minute)

				if !isWithinBusinessHours(slotStartPT, slotEndPT) {
					continue
				}

				if slotStartPT.Before(startPT) {
					continue
				}

				if slotEndPT.After(endPT) {
					continue
				}

				if anyOverlaps(slotStartPT, slotEndPT, existing) {
					continue
				}

				slots = append(slots, AvailabilitySlot{
					StartTime: slotStartPT,
					EndTime:   slotEndPT,
				})
			}
		}
	}
	return slots, nil
}

// Helper: returns the start-of-day time for a given time
func dayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Helper: returns true if any of the appointments overlaps with the given time range
func anyOverlaps(start, end time.Time, appointments []Appointment) bool {
	for _, a := range appointments {
		if start.Before(a.EndsAt) && end.After(a.StartsAt) {
			return true
		}
	}
	return false
}

// Helper: validates an appointment
func validateAppointment(appt Appointment) error {
	if appt.TrainerID == 0 || appt.UserID == 0 {
		return ErrInvalidAppointment
	}

	if !appt.EndsAt.Equal(appt.StartsAt.Add(30 * time.Minute)) {
		return ErrInvalidAppointment
	}

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}

	startPT := appt.StartsAt.In(loc)
	endPT := appt.EndsAt.In(loc)

	if !isValidStartMinute(startPT) {
		return ErrInvalidAppointment
	}
	if !isBusinessDay(startPT) {
		return ErrInvalidAppointment
	}

	if !isWithinBusinessHours(startPT, endPT) {
		return ErrInvalidAppointment
	}

	return nil
}

// Helper: returns true if the given time is a valid start minute
func isValidStartMinute(t time.Time) bool {
	return t.Minute() == 0 || t.Minute() == 30
}

// Helper: returns true if the given time is a business day
func isBusinessDay(t time.Time) bool {
	return t.Weekday() >= time.Monday && t.Weekday() <= time.Friday
}

// Helper: returns true if the given time range is within business hours
func isWithinBusinessHours(start, end time.Time) bool {
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()

	return startMinutes >= 8*60 && endMinutes <= 17*60
}
