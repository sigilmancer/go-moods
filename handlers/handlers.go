package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"go-mood-tracker/models"
	"strconv"
)

// CallingResponse belongs here now because it's part of the Handler's "contract"
type CallingResponse struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type MoodStore interface {
	Add(models.MoodEntry)
	GetAll() []models.MoodEntry
	Delete(int64)
}

type Handler struct {
	Store MoodStore
}

// Hello is a simple "health check" style handler
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from your Go server!")
}

// GetCalling handles the JSON response for Svelte component
func (h *Handler) GetCalling(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	resp := CallingResponse{
		Message:   "Hello from a Go func 🎶 - called by Svelte component",
		Timestamp: time.Now().Unix(),
	}
	
	json.NewEncoder(w).Encode(resp)
}

//single-purpose function
func (h *Handler) SubmitMood(w http.ResponseWriter, r *http.Request) {
	var entry models.MoodEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if entry.Mood == "" {
		http.Error(w, "mood is required", http.StatusBadRequest)
		return
	}

	h.Store.Add(entry)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetPulse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// If GetAll() returns nil, Go's json encoder returns 'null'. 	
	data := h.Store.GetAll()
	if data == nil {
		data = []models.MoodEntry{}
	}
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) DeleteMood(w http.ResponseWriter, r *http.Request) {
    //Get 'time' from the URL query string
    timeStr := r.URL.Query().Get("time")
    
    //Convert it to a number (int64)
    timestamp, err := strconv.ParseInt(timeStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid timestamp", http.StatusBadRequest)
        return
    }

    //Call MemoryStore to delete
    h.Store.Delete(timestamp)

    //Send back a 'Success' (204 No Content is standard for Delete)
    w.WriteHeader(http.StatusNoContent)
}
