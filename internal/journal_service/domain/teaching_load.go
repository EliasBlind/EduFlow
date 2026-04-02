package domain

// Все ID uuid7

type CreateTeachingLoad struct {
	TeacherID string
	SubjectID string
	ClassID   string
}

type TeachingLoad struct {
	ID        string
	TeacherID string
	SubjectID string
	ClassID   string
}

type UpdateTeachingLoad struct {
	ID        string
	TeacherID string
}

type ListTeachingLoadRequest struct {
	// Может быть только одно поле (oneof)
	TeacherID *string
	ClassID   *string
}

type ListTeachingLoadResponse struct {
	TotalCount uint32
	TeachingLoads []TeachingLoad
}