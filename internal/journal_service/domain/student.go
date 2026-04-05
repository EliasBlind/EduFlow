package domain

// Все ID uuid7

type Student struct {
	ID       ID
	ClassID  *ID
	FullName string
}

type CreateStudent struct {
	ClassID  *ID
	FullName string
}

type ListStudentsRequest struct {
	ClassID ID
	Limit   uint32
	Offset  uint32
}

type ListStudentsResponse struct {
	TotalCount uint32
	Students   []Student
}

type UpdateStudent struct {
	ID       ID
	ClassID  *ID
	FullName *string
}
