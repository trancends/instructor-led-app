package model

import (
	"time"
)

type Schedule struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	Date          string     `json:"date"`
	StartTime     string     `json:"start_time"`
	EndTime       string     `json:"end_time"`
	Documentation string     `json:"documentation"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
