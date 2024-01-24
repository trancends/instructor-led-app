package config

const (
	InsertIntoUsers = `INSERT INTO users (name,email,password,role,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)`
)
