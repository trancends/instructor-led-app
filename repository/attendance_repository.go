package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"enigmaCamp.com/instructor_led/config"
	"enigmaCamp.com/instructor_led/model"
)

type AttendanceRepository interface {
	GetAttendance(id string) (model.Attendance, error)
	List() ([]model.Attendance, error)
	Create(user_id string, schedule_id string) (model.Attendance, error)
	GetByID(user_id string, schedule_id string) (model.Attendance, error)
	DeleteSoft(user_id string) error
}

type attendanceRepository struct {
	db *sql.DB
}

// DeleteSoft implements AttendanceRepository.
func (a *attendanceRepository) DeleteSoft(user_id string) error {
	deleted_at := time.Now().Local()
	result, err := a.db.Exec(config.DeleteAttendance, deleted_at, user_id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetByID implements AttendanceRepository.
func (a *attendanceRepository) GetByID(user_id string, schedule_id string) (model.Attendance, error) {
	var attendance model.Attendance
	row := a.db.QueryRow(`
		SELECT id, user_id, schedule_id, created_at, updated_at, deleted_at
		FROM attendances
		WHERE user_id = $1 AND schedule_id = $2
	`, user_id, schedule_id)

	// Scan the result into the attendance struct
	err := row.Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.ScheduleID,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
		&attendance.DeletedAt,
	)
	if err != nil {
		log.Println(err)
		return model.Attendance{}, err
	}

	return attendance, nil
}

// GetAttendance implements AttendanceRepository.
func (a *attendanceRepository) GetAttendance(id string) (model.Attendance, error) {
	var attendance model.Attendance

	// Execute the SQL query
	row := a.db.QueryRow(`
		SELECT id, user_id, schedule_id
		FROM attendances
		WHERE id = $1
	`, id)

	// Scan the result into the attendance struct
	err := row.Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.ScheduleID,
	)
	if err != nil {
		log.Println(err)
		return model.Attendance{}, err
	}

	return attendance, nil
}

func (a *attendanceRepository) Create(user_id string, schedule_id string) (model.Attendance, error) {
	var attendance model.Attendance

	// Validate UUIDs

	attendance.UserID = user_id
	attendance.ScheduleID = schedule_id
	err := a.db.QueryRow(`
		INSERT INTO attendances (user_id, schedule_id)
		VALUES ($1, $2)
		RETURNING id
	`, user_id, schedule_id).Scan(
		&attendance.ID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Attendance{}, fmt.Errorf("no rows returned, record may not be inserted")
		}

		log.Printf("Error inserting attendance: %v", err)
		return model.Attendance{}, fmt.Errorf("failed to insert attendance: %v", err)
	}

	return attendance, nil
}

func (a *attendanceRepository) List() ([]model.Attendance, error) {
	var attendances []model.Attendance

	// Execute the SQL query
	rows, err := a.db.Query(`
		SELECT id, user_id, schedule_id
		FROM attendances
	`)
	if err != nil {
		log.Printf("Error querying attendances: %v", err)
		return nil, fmt.Errorf("failed to query attendances: %v", err)
	}
	defer rows.Close()

	// Iterate through the result set and populate the attendances slice
	for rows.Next() {
		var attendance model.Attendance

		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.ScheduleID,
		)
		if err != nil {
			log.Printf("Error scanning attendance row: %v", err)
			return nil, fmt.Errorf("failed to scan attendance row: %v", err)
		}

		attendances = append(attendances, attendance)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over attendance rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over attendance rows: %v", err)
	}

	return attendances, nil
}

func NewAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}
