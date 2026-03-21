package appointments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) Availability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("availability"))
}

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

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
