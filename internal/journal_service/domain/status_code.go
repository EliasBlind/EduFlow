package domain

// Все ID uuid7

type StatusCode struct {
	ID       ID
	FullName string
}

type UpdateStatusCode struct {
	ID       ID
	FullName string
}

type ListStatusCode struct {
	TotalCount  uint32
	StatusCodes []StatusCode
}
