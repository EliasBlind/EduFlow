package domain

import "time"

// Все ID uuid7

type RecordHomework struct {
	TeacherID       ID
	ClassID         ID
	SubjectID       ID
	DescriptionTask string

	// Данные ниже могут быть nil (optional)
	Start *time.Time
	End   *time.Time
}

type Homework struct {
	ID              ID
	ClassID         ID
	TeacherID       ID
	SubjectID       ID
	DescriptionTask string

	// Данные ниже могут быть nil (optional)
	Start *time.Time
	End   *time.Time
}

type UpdateHomework struct {
	ID              ID
	DescriptionTask string

	// Данные ниже могут быть nil (optional)
	Start *time.Time
	End   *time.Time
}

type ListHomeworkRequest struct {
	ClassID   ID
	SubjectID ID

	// Данные ниже могут быть nil (optional)
	Start *time.Time
	End   *time.Time
}

type ListHomeworkResponse struct {
	TotalCount uint32
	Homeworks  []Homework
}
