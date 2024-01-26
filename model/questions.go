package model

import "time"

type Question struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id,omitempty"`
	ScheduleID  string     `json:"schedule_id,omitempty"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
