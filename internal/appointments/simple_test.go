package appointments

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateAppointment_ValidHalfHourBusinessSlot(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	start := time.Date(2024, time.January, 8, 9, 0, 0, 0, loc) // Monday
	appt := Appointment{
		TrainerID: 1,
		UserID:    1,
		StartsAt:  start,
		EndsAt:    start.Add(30 * time.Minute),
	}

	if err := validateAppointment(appt); err != nil {
		t.Fatalf("expected valid appointment, got error: %v", err)
	}
}

func TestAvailability_MethodNotAllowed(t *testing.T) {
	h := NewHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/availability", nil)
	rr := httptest.NewRecorder()

	h.Availability(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	if !errors.Is(nil, nil) {
		// no-op: keeps errors import used for this tiny test file
	}
}

