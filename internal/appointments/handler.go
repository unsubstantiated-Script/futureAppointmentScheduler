package appointments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Returns all the appointments for a given trainer
func (h *Handler) getAppointments(w http.ResponseWriter, r *http.Request) {
	trainerIDstr := r.URL.Query().Get("trainer_id")
	if trainerIDstr == "" {
		http.Error(w, "trainer_id is required", http.StatusBadRequest)
		return
	}

	trainerID, err := strconv.Atoi(trainerIDstr)
	if err != nil {
		http.Error(w, "invalid trainer_id", http.StatusBadRequest)
		return
	}

	appointments, err := h.service.GetByTrainer(trainerID)
	if err != nil {
		http.Error(w, "failed to get appointments", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, appointments)
}

// Creates a new appointment
func (h *Handler) createAppointment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var appt Appointment
	if err := json.NewDecoder(r.Body).Decode(&appt); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateAppointment(appt)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAppointment):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrAppointmentOverlap):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "failed to create appointment", http.StatusInternalServerError)
		}
		return
	}

	appt.ID = id

	writeJSON(w, http.StatusCreated, appt)
}

func (h *Handler) trainerAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trainerIDstr := r.URL.Query().Get("trainer_id")
	startsAtstr := r.URL.Query().Get("starts_at")
	endsAtstr := r.URL.Query().Get("ends_at")

	if trainerIDstr == "" || startsAtstr == "" || endsAtstr == "" {
		http.Error(w, "trainer_id, starts_at, and ends_at are required", http.StatusBadRequest)
		return
	}

	trainerID, err := strconv.Atoi(trainerIDstr)
	if err != nil {
		http.Error(w, "invalid trainer_id", http.StatusBadRequest)
		return
	}

	startsAt, err := time.Parse(time.RFC3339, startsAtstr)
	if err != nil {
		http.Error(w, "invalid starts_at", http.StatusBadRequest)
		return
	}

	endsAt, err := time.Parse(time.RFC3339, endsAtstr)
	if err != nil {
		http.Error(w, "invalid ends_at", http.StatusBadRequest)
		return
	}

	slots, err := h.service.GetAvailability(trainerID, startsAt, endsAt)
	if err != nil {
		if errors.Is(err, ErrInvalidAppointment) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "failed to get trainer availability", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, slots)
}

// Appointments handles GET(list) and POST(create) on /appointments.
func (h *Handler) Appointments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAppointments(w, r)
	case http.MethodPost:
		h.createAppointment(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Availability handles GET on /availability.
func (h *Handler) Availability(w http.ResponseWriter, r *http.Request) {
	h.trainerAvailability(w, r)
}

// Helper: writes JSON to the response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
