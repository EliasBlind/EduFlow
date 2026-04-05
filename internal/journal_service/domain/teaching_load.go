package domain

// Все ID uuid7

type CreateTeachingLoad struct {
	TeacherID ID
	SubjectID ID
	ClassID   ID
}

type TeachingLoad struct {
	ID        ID
	TeacherID ID
	SubjectID ID
	ClassID   ID
}

type UpdateTeachingLoad struct {
	ID        ID
	TeacherID ID
}

type ListTeachingLoadRequest struct {
	// Может быть только одно поле (oneof)
	TeacherID *ID
	ClassID   *ID
}

type ListTeachingLoadResponse struct {
	TotalCount    uint32
	TeachingLoads []TeachingLoad
}
