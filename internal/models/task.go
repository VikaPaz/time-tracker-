package models

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID     uuid.UUID `json:"id" `
	Task   string    `json:"task"`
	UserID uuid.UUID `json:"user_id"`
}

type LaborTimeRequest struct {
	UserID *uuid.UUID `json:"user_id"`
	Limit  time.Time  `json:"limit,omitempty"`
	Offset time.Time  `json:"offset,omitempty"`
}

type LaborTimeResponse struct {
	UserID uuid.UUID  `json:"user_id"`
	Tasks  []TaskInfo `json:"tasks"`
	Total  int64      `json:"total"`
}

type TaskInfo struct {
	ID        uuid.UUID      `json:"id"`
	Task      *string        `json:"task"`
	LaborTime *time.Duration `json:"time"`
}
