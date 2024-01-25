package repository

import (
	"database/sql"
	"fmt"
	"log"

	"enigmaCamp.com/instructor_led/model"
)

type AttendanceRepository interface {
	Get(id string) (model.Attendance, error)
	List() ([]model.Attendance, error)
	Post(user_id string, schedule_id string) (model.Attendance, error)
}

type attendanceRepository struct {
	db *sql.DB
}

// Get implements AttendanceRepository.
func (*attendanceRepository) Get(id string) (model.Attendance, error) {
	panic("unimplemented")
}

// Post implements AttendanceRepository.
func (a *attendanceRepository) Post(user_id string, schedule_id string) (model.Attendance, error) {
	var attendance model.Attendance

	// Execute the SQL query to insert a new attendance record
	err := a.db.QueryRow(`
		INSERT INTO attendances (user_id, schedule_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at, deleted_at
	`, user_id, schedule_id).Scan(
		&attendance.ID,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
		&attendance.DeletedAt,
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

// List implements AttendanceRepository.
func (a *attendanceRepository) List() ([]model.Attendance, error) {
	var attendances []model.Attendance

	// Execute the SQL query
	rows, err := a.db.Query(`
		SELECT id, user_id, schedule_id, created_at, updated_at, deleted_at
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
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
			&attendance.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning attendance row: %v", err)
			return nil, fmt.Errorf("failed to scan attendance row: %v", err)
		}

		attendances = append(attendances, attendance)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over attendance rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over attendance rows: %v", err)
	}

	return attendances, nil
}
func NewAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}
