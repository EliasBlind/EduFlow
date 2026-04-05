package domain

// Все ID uuid7

type Teacher struct {
	ID       ID
	FullName string
}

type ListTeachersRequest struct {
	Limit  uint32
	Offset uint32
}

type ListTeachersResponse struct {
	TotalCount uint32
	Teachers   []Teacher
}

type UpdateTeacher struct {
	ID       ID
	FullName string
}
