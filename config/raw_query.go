package config

const (
	InsertUser = `INSERT INTO users (name,email,password,role,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
)
