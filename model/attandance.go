package model

import "time"

// ScheduleAssociation represents the structure for the provided fields.
type Attandance struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	ScheduleID string    `db:"schedule_id"`
	CreatedAt  time.Time `db:"created_at,omitempty"`
	UpdatedAt  time.Time `db:"updated_at,omitempty"`
	DeletedAt  time.Time `db:"deleted_at,omitempty"`
}
