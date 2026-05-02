package storage

import (
	"sync"
	"go-mood-tracker/models"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	activeMoods  = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mood_tracker_submissions_entries",
		Help: "The total number of processed mood entries",
	})
)

type MemoryStore struct {
	mu    sync.RWMutex
	pulse []models.MoodEntry
}

func NewStore() *MemoryStore {
	return &MemoryStore{
		pulse: make([]models.MoodEntry, 0),
	}
}

func (s *MemoryStore) Add(entry models.MoodEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pulse = append(s.pulse, entry)
	activeMoods.Inc() 
}

func (s *MemoryStore) GetAll() []models.MoodEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	//Return a copy to avoid data races with the caller
	res := make([]models.MoodEntry, len(s.pulse))
	copy(res, s.pulse)
	return res
}

func (s *MemoryStore) Delete(timestamp int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	//Filter out the entry with the matching timestamp
	newPulse := make([]models.MoodEntry, 0)
	for _, entry := range s.pulse {
		if entry.Time != timestamp {
			newPulse = append(newPulse, entry)
		}
	}
	s.pulse = newPulse
	activeMoods.Dec()
}
