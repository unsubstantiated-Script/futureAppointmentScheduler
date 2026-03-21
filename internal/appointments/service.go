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

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByTrainer(trainerID int) ([]Appointment, error) {
	return s.repo.GetByTrainer(trainerID)
}

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

func isValidStartMinute(t time.Time) bool {
	return t.Minute() == 0 || t.Minute() == 30
}

func isBusinessDay(t time.Time) bool {
	return t.Weekday() >= time.Monday && t.Weekday() <= time.Friday
}

func isWithinBusinessHours(start, end time.Time) bool {
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()

	return startMinutes >= 8*60 && endMinutes <= 17*60
}
