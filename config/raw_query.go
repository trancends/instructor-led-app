package config

const (
	SelectSchedulePagination = `
	SELECT id, user_id, date, start_time, end_time, documentation, created_at, updated_at
	FROM schedules
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2`
	InsertUser           = `INSERT INTO users (name,email,password,role,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	SelectUserPagination = `SELECT id, name, email, role FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2 WHERE deleted_at IS NULL`
)
