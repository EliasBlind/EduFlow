package domain

import "github.com/google/uuid"

type CreateClass struct {
	ClassName string `validate:"required"`
	// Лет обучается
	YearOfStudy uint32
	// Год окончания, если не установлен то будет по считан:
	// year_of_study + текущий год
	GraduationYear *uint32 `validate:"omitempty,year_gte_now"`
}

type Class struct {
	ID             uuid.UUID
	ClassName      string    `validate:"required"`
	YearOfStudy    uint32
	GraduationYear uint32 `validate:"omitempty,year_gte_now"`
}

type ListTeacherClasses struct {
	TeacherID uuid.UUID
	Limit     uint32
	Offset    uint32
}

type ListClasses struct {
	TotalCount uint32
	Classes    []Class
}

type UpdateClass struct {
	ID             uuid.UUID
	ClassName      *string   `validate:"omitempty,required"`
	YearOfStudy    *uint32
	GraduationYear *uint32 `validate:"omitempty,year_gte_now"`
}
