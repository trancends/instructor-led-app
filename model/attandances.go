package model

import "time"

type Attendance struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	ScheduleID string     `json:"schedule_id"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
