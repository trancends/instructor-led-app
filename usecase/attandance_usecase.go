// usecase/attendance_usecase.go
package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/repository"
	"github.com/google/uuid"
)

// AttendanceUsecase represents the usecase for attendance operations.
type AttendanceUsecase interface {
	GetAttendance(id string) (model.Attendance, error)
	ListAttendances() ([]model.Attendance, error)
	AddAttendance(user_id string, schedule_id string) (model.Attendance, error)
	DeleteAttandace(id string) error
}

type attendanceUsecase struct {
	attendanceRepo repository.AttendanceRepository
}

// DeleteAttandace implements AttendanceUsecase.
func (a *attendanceUsecase) DeleteAttandace(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid ID format")
	}

	return a.attendanceRepo.DeleteAttendance(id)
}

// AddAttendance implements AttendanceUsecase.
func (a *attendanceUsecase) AddAttendance(user_id string, schedule_id string) (model.Attendance, error) {
	existingAttendance, err := a.attendanceRepo.GetByID(user_id, schedule_id)

	if user_id == "" || schedule_id == "" {
		return model.Attendance{}, fmt.Errorf("user_id and schedule_id must be valid UUIDs")
	}
	if err == nil && existingAttendance.ID != "" {
		fmt.Println(existingAttendance)
		return model.Attendance{}, fmt.Errorf("attendance for user_id %s and schedule_id %s already exists", user_id, schedule_id)
	}
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking for existing attendance: %v", err)
		return model.Attendance{}, fmt.Errorf("failed to check for existing attendance: %v", err)
	}

	// If no duplicate is found, proceed to add the attendance
	return a.attendanceRepo.Create(user_id, schedule_id)
}

// NewAttendanceUsecase initializes a new AttendanceUsecase.
func NewAttendanceUsecase(attendanceRepo repository.AttendanceRepository) AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: attendanceRepo,
	}
}

// GetAttendance returns a single attendance based on the given ID.
func (u *attendanceUsecase) GetAttendance(id string) (model.Attendance, error) {
	attendance, err := u.attendanceRepo.GetAttendance(id)
	if err != nil {
		log.Println("AttendanceUsecase.GetAttendance:", err.Error())
		return model.Attendance{}, err
	}
	return attendance, nil
}

// ListAttendances returns a list of all attendances.
func (u *attendanceUsecase) ListAttendances() ([]model.Attendance, error) {
	attendances, err := u.attendanceRepo.List()
	if err != nil {
		log.Println("AttendanceUsecase.ListAttendances:", err.Error())
		return nil, err
	}
	return attendances, nil
}
