package config

const (
	SelectSchedulePagination = `
	SELECT id, user_id, date, start_time, end_time, documentation, created_at, updated_at
	FROM schedules
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2`
)
