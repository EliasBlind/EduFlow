package domain

import "time"

// Все ID uuid7

type RecordGrade struct {
	SubjectID   ID
	StudentID   ID
	DateOfGrade *time.Time

	// Может быть только одно поле (oneof)
	Grade        *uint32
	StatusCodeID *ID

	LessonNumber uint32
	Note         *string
}

type Grade struct {
	ID          ID
	SubjectID   ID
	StudentID   ID
	ClassID     ID
	DateOfGrade *time.Time

	// Может быть только одно поле (oneof)
	Grade        *uint32
	StatusCodeID ID

	LessonNumber *uint32
	Note         *string
}

type ListGradesRequest struct {
	SubjectID *ID

	// Может быть только одно поле (oneof)
	StudentID ID
	ClassID   ID

	Start *time.Time
	End   *time.Time // Может отсутствовать (optional)
}

type ListGradesResponse struct {
	TotalCount uint32
	Grades     []Grade
}

type UpdateGrade struct {
	GradeID ID

	// Может быть только одно поле (oneof)
	Grade        *uint32
	StatusCodeID *ID

	Note *string
}
