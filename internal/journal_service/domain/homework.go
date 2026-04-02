package domain

import "time"

// Все ID uuid7

type RecordHomework struct {
	TeacherID       string
	ClassID         string
	SubjectID       string
	DescriptionTask string

	// Данные ниже могут быть nil (optional)
	start *time.Time
	end   *time.Time
}

type Homework struct {
	ID              string
	ClassID         string
	TeacherID       string
	SubjectID       string
	DescriptionTask string

	// Данные ниже могут быть nil (optional)
	start *time.Time
	end   *time.Time
}

type UpdateHomework struct {
	ID              string
	DescriptionTask string

	// Данные ниже могут быть nil (optional)
	start *time.Time
	end   *time.Time
}

type ListHomeworkRequest struct {
	ClassID   string
	SubjectID string

	// Данные ниже могут быть nil (optional)
	start *time.Time
	end   *time.Time
}

type ListHomeworkResponse struct {
	TotalCount uint32
	homeworks  []Homework
}
