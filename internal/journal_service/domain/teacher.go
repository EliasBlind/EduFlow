package domain

// Все ID uuid7

type Teacher struct {
	ID       string
	FullName string
}

type ListTeachersRequest struct {
	Limit  uint32
	Offset uint32
}

type ListTeachersResponse struct {
	TotalCount uint32
	teachers   []Teacher
}

type UpdateTeacher struct {
	ID       string
	FullName string
}
