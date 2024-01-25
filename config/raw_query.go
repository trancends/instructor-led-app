package config

const (
	SelectSchedulePagination = `
	SELECT id, user_id, date, start_time, end_time, documentation
	FROM schedules
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2`
	InsertUser         = `INSERT INTO users (name,email,password,role) VALUES ($1,$2,$3,$4) RETURNING id`
	InsertQuestions    = `INSERT INTO questions (id,schedule_id,description) VALUES ($1,$2,$3,'PROCESSED') RETURNING id`
	InsertSchedule     = `INSERT INTO schedules (id,user_id,date,start_time,end_time,documentation) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	SelectScheduleByID = `SELECT id, user_id, documentation, date, start_time, end_time FROM schedules WHERE id = $1`
)
