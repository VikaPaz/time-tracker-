package models

import (
	"github.com/google/uuid"
	"time"
)

type UserTask struct {
	UserID *uuid.UUID `json:"user_id"`
	Text   *string    `json:"text"`
}

type Timer struct {
	TaskID uuid.UUID `json:"task_id"`
}

type Task struct {
	ID     uuid.UUID `json:"id" `
	Task   string    `json:"task"`
	UserID uuid.UUID `json:"user_id"`
}

type LaborTimeRequest struct {
	UserID    *uuid.UUID `json:"user_id"`
	StartTime *time.Time `json:"limit,omitempty"`
	EndTime   *time.Time `json:"offset,omitempty"`
}

type LaborTimeResponse struct {
	UserID uuid.UUID  `json:"user_id"`
	Tasks  []TaskInfo `json:"tasks"`
}

type TaskInfo struct {
	ID        uuid.UUID      `json:"id"`
	Task      *string        `json:"task"`
	LaborTime *time.Duration `json:"time"`
}

type GetTaskResponse struct {
	UserID uuid.UUID     `json:"user_id"`
	Tasks  []GetTaskInfo `json:"tasks"`
}

type GetTaskInfo struct {
	ID        uuid.UUID `json:"id"`
	Task      string    `json:"task"`
	LaborTime string    `json:"time"`
}
