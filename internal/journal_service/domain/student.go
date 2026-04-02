package domain

// Все ID uuid7

type Student struct {
	ID       string
	ClassID  *string
	FullName string
}

type CreateStudent struct {
	ClassID  *string
	FullName string
}

type ListStudentsRequest struct {
	ClassID string
	Limit   uint32
	Offset  uint32
}

type ListStudentsResult struct {
	TotalCount uint32
	students   []Student
}

type UpdateStudent struct {
	ID       string
	ClassID  *string
	FullName *string
}
