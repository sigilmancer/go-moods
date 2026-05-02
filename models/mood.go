package models

type MoodEntry struct {
    Time int64  `json:"time"`
    Mood string `json:"mood"`
}
