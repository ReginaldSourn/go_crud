package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// Adjust the import path according to your project structure
)

type DevicesHandler struct {
	store store.DevicesStore
}

// NewDevicesHandler constructs a new DevicesHandler.
// Replace store.DevicesStore with whatever type or interface you defined in store/devices.go.
func NewDevicesHandler(s store.DevicesStore) *DevicesHandler {
	return &DevicesHandler{store: s}
}

// CreateDevice handles POST /devices
func (h *DevicesHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var d store.Device
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.store.CreateDevice(r.Context(), &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	d.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

// GetDevice handles GET /devices/{id}
func (h *DevicesHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	d, err := h.store.GetDevice(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// UpdateDevice handles PUT /devices/{id}
func (h *DevicesHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var d store.Device
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d.ID = id
	if err := h.store.UpdateDevice(r.Context(), &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// DeleteDevice handles DELETE /devices/{id}
func (h *DevicesHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteDevice(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListDevices handles GET /devices
func (h *DevicesHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	ds, err := h.store.ListDevices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ds)
}
