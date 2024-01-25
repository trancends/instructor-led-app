package model

import "time"

type Attendance struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	ScheduleID string     `json:"schedule_id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
