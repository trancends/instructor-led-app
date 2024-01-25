package model

import "time"

type Questions struct {
	ID          string     `json:"id"`
	ScheduleID  string     `json:"schedule_id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
