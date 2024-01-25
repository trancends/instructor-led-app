package config

const (
	SelectSchedulePagination = `
	SELECT id, user_id, date, start_time, end_time, documentation
	FROM schedules
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2`

	// User
	InsertUser           = `INSERT INTO users (name,email,password,role) VALUES ($1,$2,$3,$4) RETURNING id`
	SelectUserPagination = `SELECT id, name, email, role FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2 WHERE deleted_at IS NULL`
	SelectUserByID       = `SELECT id, name, email, role FROM users WHERE id = $1 AND deleted_at IS NULL`
	SelectUserByEmail    = `SELECT id, name, email, role FROM users WHERE email = $1 AND deleted_at IS NULL`
	SelectUserByRole     = `SELECT id, name, email, role FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2 WHERE role = $3 AND deleted_at IS NULL`
	UpdateUser           = `UPDATE users SET name = $1, email = $2, password = $3, role = $4, updated_at = $5 WHERE id = $6`
	DeleteUser           = `UPDATE users SET deleted_at = $1 WHERE id = $2`
	InsertQuestions      = `INSERT INTO questions (id,schedule_id,description) VALUES ($1,$2,$3,'PROCESSED') RETURNING id`
	InsertSchedule       = `INSERT INTO schedules (user_id, date, start_time, end_time, documentation) VALUES ($1, $2, $3, $4, $5) RETURNING id,user_id, date, start_time, end_time, documentation`
	SelectScheduleByID   = `SELECT id, user_id, date, start_time, end_time, documentation FROM schedules WHERE id = $1`
	DeleteSchedule       = `UPDATE schedules SET deleted_at = $1 WHERE id = $2`
)
