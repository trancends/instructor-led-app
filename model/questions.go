package model

type Question struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id,omitempty"`
	ScheduleID  string `json:"schedule_id,omitempty"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	DeletedAt   string `json:"deleted_at,omitempty"`
}
