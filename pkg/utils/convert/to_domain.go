package convert

import (
	"time"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
)

// Маппинг параметров функций (toDomain)

// Студенты (Students)
func CreateStudentToDomain(j *journalv1.CreateStudentRequest) *domain.CreateStudent {
	return &domain.CreateStudent{
		// option field
		ClassID:  (*domain.ID)(j.ClassId),
		FullName: j.GetFullName(),
	}
}

func ListStudentsToDomain(j *journalv1.ListStudentsRequest) *domain.ListStudentsRequest {
	return &domain.ListStudentsRequest{
		ClassID: domain.ID(j.ClassId),
		Limit:   j.GetLimit(),
		Offset:  j.GetOffset(),
	}
}

func UpdateStudentToDomain(j *journalv1.UpdateStudentRequest) *domain.UpdateStudent {
	return &domain.UpdateStudent{
		ID: domain.ID(j.Id),

		// option fields
		ClassID:  (*domain.ID)(j.ClassId),
		FullName: j.FullName,
	}
}

// Учителя (Teachers)
func ListTeachersToDomain(j *journalv1.ListTeachersRequest) *domain.ListTeachersRequest {
	return &domain.ListTeachersRequest{
		Limit:  j.GetLimit(),
		Offset: j.GetOffset(),
	}
}

func UpdateTeacherToDomain(j *journalv1.UpdateTeacherRequest) *domain.UpdateTeacher {
	return &domain.UpdateTeacher{
		ID:       domain.ID(j.Id),
		FullName: j.GetFullName(),
	}
}

// Классы (Classes)
func CreateClassToDomain(j *journalv1.CreateClassRequest) *domain.CreateClass {
	return &domain.CreateClass{
		ClassName:   j.GetClassName(),
		YearOfStudy: j.GetYearOfStudy(),

		// option field
		GraduationYear: j.GraduationYear,
	}
}

func ListTeacherClassesToDomain(j *journalv1.ListTeacherClassesRequest) *domain.ListTeacherClasses {
	return &domain.ListTeacherClasses{
		TeacherID: domain.ID(j.TeacherId),
		Limit:     j.GetLimit(),
		Offset:    j.GetOffset(),
	}
}

func UpdateClassToDomain(j *journalv1.UpdateClassRequest) *domain.UpdateClass {
	return &domain.UpdateClass{
		ID: domain.ID(j.Id),

		// option fields
		ClassName:      j.ClassName,
		YearOfStudy:    j.YearOfStudy,
		GraduationYear: j.GraduationYear,
	}
}

// Предметы (Subjects)
func UpdateSubjectToDomain(j *journalv1.UpdateSubjectRequest) *domain.UpdateSubject {
	return &domain.UpdateSubject{
		ID:       domain.ID(j.GetId()),
		FullName: j.GetFullName(),
	}
}

// Оценки (Grades)
func RecordGradeToDomain(j *journalv1.RecordGradeRequest) *domain.RecordGrade {

	var grade *uint32
	var statusCodeID *domain.ID

	switch v := j.Filter.(type) {
	case *journalv1.RecordGradeRequest_Grade:
		grade = ptr(v.Grade)

	case *journalv1.RecordGradeRequest_StatusCodeId:
		statusCodeID = ptr(domain.ID(v.StatusCodeId))
	}

	var dateOfGrade *time.Time

	if ts := j.GetDateOfGrade(); ts != nil {
		dateOfGrade = ptr(ts.AsTime())
	}

	return &domain.RecordGrade{
		SubjectID: domain.ID(j.SubjectId),
		StudentID: domain.ID(j.StudentId),

		// option field
		DateOfGrade: dateOfGrade,

		// oneof
		Grade:        grade,
		StatusCodeID: statusCodeID,

		LessonNumber: j.GetLessonNumber(),

		// option field
		Note: j.Note,
	}
}

func ListGradesToDomain(j *journalv1.ListGradesRequest) *domain.ListGradesRequest {
	start := j.GetStart().AsTime()
	end := j.GetEnd().AsTime()
	return &domain.ListGradesRequest{
		SubjectID: (*domain.ID)(j.SubjectId),

		// oneof
		StudentID: domain.ID(j.GetStudentId()),
		ClassID:   domain.ID(j.GetClassId()),

		// options fields
		Start: &start,
		End:   &end,
	}
}

func UpdateGradeToDomain(j *journalv1.UpdateGradeRequest) *domain.UpdateGrade {
	var grade *uint32
	var statusCodeID *string

	switch v := j.Filter.(type) {
	case *journalv1.UpdateGradeRequest_Grade:
		grade = &v.Grade
	case *journalv1.UpdateGradeRequest_StatusCodeId:
		statusCodeID = &v.StatusCodeId
	}

	return &domain.UpdateGrade{
		GradeID: domain.ID(j.GradeId),

		// one of
		Grade:        grade,
		StatusCodeID: (*domain.ID)(statusCodeID),

		// option field
		Note: j.Note,
	}
}

func RecordHomeworkToDomain(j *journalv1.RecordHomeworkRequest) *domain.RecordHomework {
	var start, end *time.Time

	if ts := j.GetStart(); ts != nil {
		start = ptr(ts.AsTime())
	}

	if ts := j.GetEnd(); ts != nil {
		end = ptr(ts.AsTime())
	}

	return &domain.RecordHomework{
		TeacherID:       domain.ID(j.TeacherId),
		ClassID:         domain.ID(j.ClassId),
		SubjectID:       domain.ID(j.SubjectId),
		DescriptionTask: j.GetDescriptionTask(),

		// options fields
		Start: start,
		End:   end,
	}
}

func UpdateHomeworkToDomain(j *journalv1.UpdateHomeworkRequest) *domain.UpdateHomework {
	var start, end *time.Time

	if ts := j.GetStart(); ts != nil {
		start = ptr(ts.AsTime())
	}

	if ts := j.GetEnd(); ts != nil {
		end = ptr(ts.AsTime())
	}

	return &domain.UpdateHomework{
		ID:              domain.ID(j.GetId()),
		DescriptionTask: j.GetDescriptionTask(),

		// optional fields
		Start: start,
		End:   end,
	}
}

func ListHomeworkToDomain(j *journalv1.ListHomeworkRequest) *domain.ListHomeworkRequest {
	var start, end *time.Time

	if ts := j.GetStart(); ts != nil {
		start = ptr(ts.AsTime())
	}

	if ts := j.GetEnd(); ts != nil {
		end = ptr(ts.AsTime())
	}

	return &domain.ListHomeworkRequest{
		ClassID:   domain.ID(j.GetClassId()),
		SubjectID: domain.ID(j.GetSubjectId()),

		// options fields
		Start: start,
		End:   end,
	}
}

// Статус-коды (StatusCode)
func UpdateStatusCodeToDomain(j *journalv1.UpdateStatusCodeRequest) *domain.UpdateStatusCode {
	return &domain.UpdateStatusCode{
		ID:       domain.ID(j.GetId()),
		FullName: j.GetFullName(),
	}
}

// Нагрузка на учителей (TeachingLoad)
func CreateTeachingLoadToDomain(j *journalv1.CreateTeachingLoadRequest) *domain.CreateTeachingLoad {
	return &domain.CreateTeachingLoad{
		TeacherID: domain.ID(j.GetTeacherId()),
		SubjectID: domain.ID(j.GetSubjectId()),
		ClassID:   domain.ID(j.GetClassId()),
	}
}

func UpdateTeachingLoadToDomain(j *journalv1.UpdateTeachingLoadRequest) *domain.UpdateTeachingLoad {
	return &domain.UpdateTeachingLoad{
		ID:        domain.ID(j.GetId()),
		TeacherID: domain.ID(j.GetTeacherId()),
	}
}

func ListTeachingLoadToDomain(j *journalv1.ListTeachingLoadRequest) *domain.ListTeachingLoadRequest {
	var teacherID *domain.ID
	var classID *domain.ID

	switch v := j.Filter.(type) {
	case *journalv1.ListTeachingLoadRequest_TeacherId:
		id := domain.ID(v.TeacherId)
		teacherID = &id
	case *journalv1.ListTeachingLoadRequest_ClassId:
		id := domain.ID(v.ClassId)
		classID = &id
	}

	return &domain.ListTeachingLoadRequest{
		TeacherID: teacherID,
		ClassID:   classID,
	}
}

// Support functions

func ptr[T any](a T) *T {
	return &a
}

func IDToProtoPtr(id *domain.ID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}
