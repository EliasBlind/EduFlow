package convert

import (
	"time"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
)

func ptr[T any](a T) *T {
	return &a
}

// Маппинг параметров функций (toDomain)

// Студенты (Students)
func CreateStudentToDomain(j *journalv1.CreateStudentRequest) *domain.CreateStudent {
	return &domain.CreateStudent{
		// option field
		ClassID:  j.ClassId,
		FullName: j.GetFullName(),
	}
}

func ListStudentsToDomain(j *journalv1.ListStudentsRequest) *domain.ListStudentsRequest {
	return &domain.ListStudentsRequest{
		ClassID: j.ClassId,
		Limit:   j.GetLimit(),
		Offset:  j.GetOffset(),
	}
}

func UpdateStudentToDomain(j *journalv1.UpdateStudentRequest) *domain.UpdateStudent {
	return &domain.UpdateStudent{
		ID: j.Id,

		// option fields
		ClassID:  j.ClassId,
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
		ID:       j.Id,
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
		TeacherID: j.TeacherId,
		Limit:     j.GetLimit(),
		Offset:    j.GetOffset(),
	}
}

func UpdateClassToDomain(j *journalv1.UpdateClassRequest) *domain.UpdateClass {
	return &domain.UpdateClass{
		ID: j.Id,

		// option fields
		ClassName:      j.ClassName,
		YearOfStudy:    j.YearOfStudy,
		GraduationYear: j.GraduationYear,
	}
}

// Предметы (Subjects)
func UpdateSubjectToDomain(j *journalv1.UpdateSubjectRequest) *domain.UpdateSubject {
	return &domain.UpdateSubject{
		ID:       j.GetId(),
		FullName: j.GetFullName(),
	}
}

// Оценки (Grades)
func RecordGradeToDomain(j *journalv1.RecordGradeRequest) *domain.RecordGrade {

	var grade *uint32
	var statusCodeID *string

	switch v := j.Filter.(type) {
	case *journalv1.RecordGradeRequest_Grade:
		grade = ptr(v.Grade)

	case *journalv1.RecordGradeRequest_StatusCodeId:
		statusCodeID = ptr(v.StatusCodeId)
	}

	var dateOfGrade *time.Time

	if ts := j.GetDateOfGrade(); ts != nil {
		dateOfGrade = ptr(ts.AsTime())
	}

	return &domain.RecordGrade{
		SubjectID: j.SubjectId,
		StudentID: j.StudentId,

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
		SubjectID: j.SubjectId,

		// oneof
		StudentID: j.GetStudentId(),
		ClassID:   j.GetClassId(),

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
		GradeID: j.GradeId,

		// one of
		Grade:        grade,
		StatusCodeID: statusCodeID,

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
		TeacherID:       j.TeacherId,
		ClassID:         j.ClassId,
		SubjectID:       j.SubjectId,
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
		ID:              j.Id,
		DescriptionTask: j.GetDescriptionTask(),

		// options fields
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
		ClassID:   j.GetClassId(),
		SubjectID: j.GetSubjectId(),

		// options fields
		Start: start,
		End:   end,
	}
}

// Статус-коды (StatusCode)
func UpdateStatusCodeToDomain(j *journalv1.UpdateStatusCodeRequest) *domain.UpdateStatusCode {
	return &domain.UpdateStatusCode{
		ID:       j.GetId(),
		FullName: j.GetFullName(),
	}
}

// Нагрузка на учителей (TeachingLoad)
func CreateTeachingLoadToDomain(j *journalv1.CreateTeachingLoadRequest) *domain.CreateTeachingLoad {
	return &domain.CreateTeachingLoad{
		TeacherID: j.GetTeacherId(),
		SubjectID: j.GetSubjectId(),
		ClassID:   j.GetClassId(),
	}
}

func UpdateTeachingLoadToDomain(j *journalv1.UpdateTeachingLoadRequest) *domain.UpdateTeachingLoad {
	return &domain.UpdateTeachingLoad{
		ID:        j.GetId(),
		TeacherID: j.GetTeacherId(),
	}
}

func ListTeachingLoadToDomain(j *journalv1.ListTeachingLoadRequest) *domain.ListTeachingLoadRequest {
	teacher_id := j.GetTeacherId()
	class_id := j.GetClassId()

	return &domain.ListTeachingLoadRequest{
		// one of fields
		TeacherID: &teacher_id,
		ClassID:   &class_id,
	}
}
