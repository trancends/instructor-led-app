// usecase/attendance_usecase.go
package usecase

import (
	"log"

	"enigmaCamp.com/instructor_led/model"
	"enigmaCamp.com/instructor_led/repository"
)

// AttendanceUsecase represents the usecase for attendance operations.
type AttendanceUsecase interface {
	GetAttendance(id string) (model.Attendance, error)
	ListAttendances() ([]model.Attendance, error)
	AddAttendance(user_id string, schedule_id string) (model.Attendance, error)
}

type attendanceUsecase struct {
	attendanceRepo repository.AttendanceRepository
}

// AddAttendance implements AttendanceUsecase.
func (a *attendanceUsecase) AddAttendance(user_id string, schedule_id string) (model.Attendance, error) {
	return a.attendanceRepo.Post(user_id, schedule_id)
}

// NewAttendanceUsecase initializes a new AttendanceUsecase.
func NewAttendanceUsecase(attendanceRepo repository.AttendanceRepository) AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: attendanceRepo,
	}
}

// GetAttendance returns a single attendance based on the given ID.
func (u *attendanceUsecase) GetAttendance(id string) (model.Attendance, error) {
	attendance, err := u.attendanceRepo.Get(id)
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
