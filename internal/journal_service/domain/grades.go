package domain

import "time"

// Все ID uuid7

type RecordGrade struct {
	SubjectID   string
	StudentID   string
	DateOfGrade *time.Time

	// Может быть только одно поле (oneof)
	Grade        *uint32
	StatusCodeID *string

	LessonNumber uint32
	Note         *string
}

type Grade struct {
	ID          string
	SubjectID   string
	StudentID   string
	ClassID     string
	DateOfGrade *time.Time

	// Может быть только одно поле (oneof)
	Grade        *uint32
	StatusCodeID string

	LessonNumber *uint32
	Note         *string
}

type ListGradesRequest struct {
	SubjectID *string

	// Может быть только одно поле (oneof)
	StudentID string
	ClassID   string

	start *time.Time
	end   *time.Time // Может отсутствовать (optional)
}

type ListGradesResponse struct {
	TotalCount uint32
	Grades     []Grade
}

type UpdateGrade struct {
	GradeID string

	// Может быть только одно поле (oneof)
	Grade        uint16
	StatusCodeID string

	Note *string
}
