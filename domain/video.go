package domain

import "time"

type Video struct {
	ID        string
	ResouceID string
	FilePath  string
	CreatedAt time.Time
}
