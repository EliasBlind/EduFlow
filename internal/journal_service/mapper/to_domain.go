package mapper

import (
	"time"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	pkgmapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"
)

// Маппинг параметров функций (toDomain)

// Студенты (Students)
func CreateStudentToDomain(j *journalv1.CreateStudentRequest) (*domain.Student, error) {

	id, err := pkgmapper.ToUUID(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	classID, err := pkgmapper.ToUUID(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.Student{
		// option field
		ID: id,
		// option field
		ClassID:  classID,
		FullName: j.GetFullName(),
	}, nil
}

func ListStudentsToDomain(j *journalv1.ListStudentsRequest) (*domain.ListStudentsRequest, error) {
	classID, err := uuid.Parse(j.ClassId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &domain.ListStudentsRequest{
		ClassID: classID,
		Limit:   j.GetLimit(),
		Offset:  j.GetOffset(),
	}, nil
}

func UpdateStudentToDomain(j *journalv1.UpdateStudentRequest) (*domain.UpdateStudent, error) {
	id, err := uuid.Parse(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	classID, err := pkgmapper.ToUUID(j.ClassId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.UpdateStudent{
		ID: id,

		// option fields
		ClassID:  classID,
		FullName: j.FullName,
	}, nil
}

// Учителя (Teachers)
func TeacherRequestToDomain(j *journalv1.CreateTeacherRequest) (*domain.Teacher, error) {
	id, err := pkgmapper.ToUUID(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.Teacher{
		ID:       id,
		FullName: j.FullName,
	}, nil
}

func ListTeachersToDomain(j *journalv1.ListTeachersRequest) *domain.ListTeachersRequest {
	return &domain.ListTeachersRequest{
		Limit:  j.GetLimit(),
		Offset: j.GetOffset(),
	}
}

func UpdateTeacherToDomain(j *journalv1.UpdateTeacherRequest) (*domain.UpdateTeacher, error) {
	id, err := uuid.Parse(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &domain.UpdateTeacher{
		ID:       id,
		FullName: j.GetFullName(),
	}, nil
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

func ListTeacherClassesToDomain(j *journalv1.ListTeacherClassesRequest) (*domain.ListTeacherClasses, error) {
	teacherID, err := uuid.Parse(j.TeacherId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &domain.ListTeacherClasses{
		TeacherID: teacherID,
		Limit:     j.GetLimit(),
		Offset:    j.GetOffset(),
	}, nil
}

func UpdateClassToDomain(j *journalv1.UpdateClassRequest) (*domain.UpdateClass, error) {
	id, err := uuid.Parse(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &domain.UpdateClass{
		ID: id,

		// option fields
		ClassName:      j.ClassName,
		YearOfStudy:    j.YearOfStudy,
		GraduationYear: j.GraduationYear,
	}, nil
}

// Предметы (Subjects)
func UpdateSubjectToDomain(j *journalv1.UpdateSubjectRequest) (*domain.UpdateSubject, error) {
	id, err := uuid.Parse(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &domain.UpdateSubject{
		ID:       id,
		FullName: j.GetFullName(),
	}, nil
}

// Оценки (Grades)
func RecordGradeToDomain(j *journalv1.RecordGradeRequest) (*domain.RecordGrade, error) {

	var grade *uint32
	var statusCodeID *uuid.UUID

	switch v := j.Filter.(type) {
	case *journalv1.RecordGradeRequest_Grade:
		grade = ptr(v.Grade)

	case *journalv1.RecordGradeRequest_StatusCodeId:
		sci, err := uuid.Parse(v.StatusCodeId)
		if err != nil {
			return nil, domain.ErrInvalidData
		}
		statusCodeID = ptr(sci)
	}

	var dateOfGrade *time.Time

	if ts := j.GetDateOfGrade(); ts != nil {
		dateOfGrade = ptr(ts.AsTime())
	}

	tsId, err := uuid.Parse(j.GetTsId())
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	studentId, err := uuid.Parse(j.GetStudentId())
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.RecordGrade{
		TsID:      tsId,
		StudentID: studentId,

		// option field
		DateOfGrade: dateOfGrade,

		// oneof
		Grade:        grade,
		StatusCodeID: statusCodeID,

		LessonNumber: j.GetLessonNumber(),

		// option field
		Note: j.Note,
	}, nil
}

func ListGradesToDomain(j *journalv1.ListGradesRequest) (*domain.ListGradesRequest, error) {
	var studentID, classID *uuid.UUID

	switch v := j.Filter.(type) {
	case *journalv1.ListGradesRequest_StudentId:
		sID, err := uuid.Parse(v.StudentId)
		if err != nil {
			return nil, domain.ErrInvalidData
		}
		studentID = ptr(sID)

	case *journalv1.ListGradesRequest_ClassId:
		cID, err := uuid.Parse(v.ClassId)
		if err != nil {
			return nil, domain.ErrInvalidData
		}
		classID = ptr(cID)
	}

	var start, end *time.Time

	if ts := j.GetStart(); ts != nil {
		start = ptr(ts.AsTime())
	}

	if ts := j.GetEnd(); ts != nil {
		end = ptr(ts.AsTime())

	}

	var subjectID *uuid.UUID
	if j.SubjectId != nil {
		sID, err := uuid.Parse(*j.SubjectId)
		if err != nil {
			return nil, domain.ErrInvalidData
		}
		subjectID = &sID
	}

	return &domain.ListGradesRequest{
		SubjectID: subjectID,

		// oneof
		StudentID: studentID,
		ClassID:   classID,

		// options fields
		Start: start,
		End:   end,
	}, nil
}

func UpdateGradeToDomain(j *journalv1.UpdateGradeRequest) (*domain.UpdateGrade, error) {
	var grade *uint32
	var statusCodeID *uuid.UUID

	switch v := j.Filter.(type) {
	case *journalv1.UpdateGradeRequest_Grade:
		grade = &v.Grade
	case *journalv1.UpdateGradeRequest_StatusCodeId:
		sCId, err := uuid.Parse(v.StatusCodeId)
		statusCodeID = &sCId
		if err != nil {
			return nil, domain.ErrInvalidData
		}
	}

	gradeID, err := uuid.Parse(j.GradeId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &domain.UpdateGrade{
		GradeID: gradeID,

		// one of
		Grade:        grade,
		StatusCodeID: statusCodeID,

		// option field
		Note: j.Note,
	}, nil
}

func RecordHomeworkToDomain(j *journalv1.RecordHomeworkRequest) (*domain.RecordHomework, error) {
	var start, end *time.Time

	if ts := j.GetStart(); ts != nil {
		start = ptr(ts.AsTime())
	}

	if ts := j.GetEnd(); ts != nil {
		end = ptr(ts.AsTime())
	}

	teacherId, err := uuid.Parse(j.GetTeacherId())
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	classId, err := uuid.Parse(j.GetClassId())
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	subjectID, err := uuid.Parse(j.GetSubjectId())
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.RecordHomework{
		TeacherID:       teacherId,
		ClassID:         classId,
		SubjectID:       subjectID,
		DescriptionTask: j.GetDescriptionTask(),

		// options fields
		Start: start,
		End:   end,
	}, nil
}

func UpdateHomeworkToDomain(j *journalv1.UpdateHomeworkRequest) (*domain.UpdateHomework, error) {
	var start, end *time.Time

	if ts := j.GetStart(); ts != nil {
		start = ptr(ts.AsTime())
	}

	if ts := j.GetEnd(); ts != nil {
		end = ptr(ts.AsTime())
	}

	id, err := uuid.Parse(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.UpdateHomework{
		ID:              id,
		DescriptionTask: j.DescriptionTask,

		// optional fields
		Start: start,
		End:   end,
	}, nil
}

func ListHomeworkToDomain(j *journalv1.ListHomeworkRequest) (*domain.ListHomeworkRequest, error) {
	var start, end *time.Time

	if ts := j.GetStart(); ts != nil {
		start = ptr(ts.AsTime())
	}

	if ts := j.GetEnd(); ts != nil {
		end = ptr(ts.AsTime())
	}

	classID, err := uuid.Parse(j.ClassId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	subjectID, err := uuid.Parse(j.SubjectId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.ListHomeworkRequest{
		ClassID:   classID,
		SubjectID: subjectID,

		// options fields
		Start: start,
		End:   end,
	}, nil
}

// Статус-коды (StatusCode)
func UpdateStatusCodeToDomain(j *journalv1.UpdateStatusCodeRequest) (*domain.UpdateStatusCode, error) {
	id, err := uuid.Parse(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &domain.UpdateStatusCode{
		ID:       id,
		FullName: j.GetFullName(),
	}, nil
}

// Нагрузка на учителей (TeachingLoad)
func CreateTeachingLoadToDomain(j *journalv1.CreateTeachingLoadRequest) (*domain.CreateTeachingLoad, error) {
	teacherID, err := uuid.Parse(j.TeacherId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	subjectID, err := uuid.Parse(j.SubjectId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	classID, err := uuid.Parse(j.ClassId)

	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.CreateTeachingLoad{
		TeacherID: teacherID,
		SubjectID: subjectID,
		ClassID:   classID,
	}, nil
}

func UpdateTeachingLoadToDomain(j *journalv1.UpdateTeachingLoadRequest) (*domain.UpdateTeachingLoad, error) {
	id, err := uuid.Parse(j.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	teacherID, err := uuid.Parse(j.TeacherId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.UpdateTeachingLoad{
		ID:        id,
		TeacherID: teacherID,
	}, nil
}

func ListTeachingLoadToDomain(j *journalv1.ListTeachingLoadRequest) (*domain.ListTeachingLoadRequest, error) {
	var teacherID, classID *uuid.UUID

	switch v := j.Filter.(type) {
	case *journalv1.ListTeachingLoadRequest_TeacherId:
		id, err := uuid.Parse(v.TeacherId)
		if err != nil {
			return nil, domain.ErrInvalidData
		}
		teacherID = &id
	case *journalv1.ListTeachingLoadRequest_ClassId:
		id, err := uuid.Parse(v.ClassId)
		if err != nil {
			return nil, domain.ErrInvalidData
		}
		classID = &id
	}

	return &domain.ListTeachingLoadRequest{
		TeacherID: teacherID,
		ClassID:   classID,
	}, nil
}

// Support functions

func ptr[T any](a T) *T {
	return &a
}
