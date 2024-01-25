package model

type Question struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	ScheduleID  string `json:"schedule_id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Date        []Schedule
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}
