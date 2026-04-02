package domain

// Все ID uuid7

type CreateClass struct {
	ClassName string
	// Лет обучается
	YearOfStudy uint32
	// Год окончания, если не установлен то будет по считан:
	// year_of_study + текущий год
	GraduationYear *uint32
}

type Class struct {
	ID             string
	ClassName      string
	YearOfStudy    uint32
	GraduationYear uint32
}

type ListTeacherClasses struct {
	TeacherID string
	Limit     uint32
	Offset    uint32
}

type ListClasses struct {
	TotalCount uint32
	Classes    []Class
}

type UpdateClass struct {
	ID             string
	ClassName      *string
	YearOfStudy    *uint32
	GraduationYear *uint32
}
