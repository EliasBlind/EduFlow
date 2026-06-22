package mapper

import (
	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	pkgmapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Студенты (Students)

func StudentToProto(params *domain.Student) *journalv1.Student {
	return &journalv1.Student{
		Id:       params.ID.String(),
		ClassId:  uuidPtrToString(params.ClassID),
		FullName: params.FullName,
	}
}

func ListStudentsToProto(students []domain.Student) *journalv1.ListStudentsResponse {
	return &journalv1.ListStudentsResponse{
		TotalCount: uint32(len(students)),
		Students:   pkgmapper.MapSliceRef(students, StudentToProto),
	}
}

// Учителя (Teacher)

func TeacherToProto(params *domain.Teacher) *journalv1.Teacher {
	return &journalv1.Teacher{
		Id:       params.ID.String(),
		FullName: params.FullName,
	}
}

func ListTeachersToProto(teacher []domain.Teacher) *journalv1.ListTeachersResponse {
	return &journalv1.ListTeachersResponse{
		TotalCount: uint32(len(teacher)),
		Teachers:   pkgmapper.MapSliceRef(teacher, TeacherToProto),
	}
}

// Класс (Classes)

func ClassToProto(params *domain.Class) *journalv1.Class {
	return &journalv1.Class{
		Id:             params.ID.String(),
		ClassName:      params.ClassName,
		YearOfStudy:    params.YearOfStudy,
		GraduationYear: params.GraduationYear,
	}
}

func ListTeacherClassesToProto(classes []domain.Class) *journalv1.ListClassesResponse {
	return &journalv1.ListClassesResponse{
		TotalCount: uint32(len(classes)),
		Classes:    pkgmapper.MapSliceRef(classes, ClassToProto),
	}
}

// Предмет (Subject)

func SubjectToProto(params *domain.Subject) *journalv1.Subject {
	return &journalv1.Subject{
		Id:       params.ID.String(),
		FullName: params.FullName,
	}
}

// Оценки (Grades)

func GradeToProto(params *domain.Grade) *journalv1.Grade {
	res := &journalv1.Grade{
		Id:        params.ID.String(),
		SubjectId: params.SubjectID.String(),
		StudentId: params.StudentID.String(),
		ClassId:   params.ClassID.String(),
		Note:      params.Note,
	}

	if params.LessonNumber != nil {
		res.LessonNumber = *params.LessonNumber
	}

	if params.DateOfGrade != nil {
		res.DateOfGrade = timestamppb.New(*params.DateOfGrade)
	}

	if params.Grade != nil {
		res.Filter = &journalv1.Grade_Grade{
			Grade: *params.Grade,
		}
	} else if params.StatusCodeID != uuid.Nil {
		res.Filter = &journalv1.Grade_StatusCodeId{
			StatusCodeId: params.StatusCodeID.String(),
		}
	}

	return res
}

func GradesToProto(grades []domain.Grade) *journalv1.ListGradesResponse {
	return &journalv1.ListGradesResponse{
		TotalCount: uint32(len(grades)),
		Grades:     pkgmapper.MapSliceRef(grades, GradeToProto),
	}
}

// Домашняя работа (Homework)

func HomeworkToProto(params *domain.Homework) *journalv1.Homework {
	return &journalv1.Homework{
		Id:              params.ID.String(),
		ClassId:         params.ClassID.String(),
		TeacherId:       params.TeacherID.String(),
		SubjectId:       params.SubjectID.String(),
		DescriptionTask: params.DescriptionTask,

		Start: timestamppb.New(*params.Start),
		End:   timestamppb.New(*params.End),
	}
}

func ListHomeworkToProto(homeworks []domain.Homework) *journalv1.ListHomeworkResponse {
	return &journalv1.ListHomeworkResponse{
		TotalCount: uint32(len(homeworks)),
		Homeworks:  pkgmapper.MapSliceRef(homeworks, HomeworkToProto),
	}
}

// Статус - коды (TeachingLoad)

func StatusCodeToProto(params *domain.StatusCode) *journalv1.StatusCode {
	return &journalv1.StatusCode{
		Id:       params.ID.String(),
		FullName: params.FullName,
	}
}

func ListStatusCodeToProto(statusCodes []domain.StatusCode) *journalv1.ListStatusCodeResponse {
	return &journalv1.ListStatusCodeResponse{
		TotalCount:  uint32(len(statusCodes)),
		StatusCodes: pkgmapper.MapSliceRef(statusCodes, StatusCodeToProto),
	}
}

// Учебная нагрузка (TeachingLoad)

func TeachingLoadToProto(params *domain.TeachingLoad) *journalv1.TeachingLoad {
	return &journalv1.TeachingLoad{
		Id:        params.ID.String(),
		TeacherId: params.TeacherID.String(),
		SubjectId: params.SubjectID.String(),
		ClassId:   params.ClassID.String(),
	}
}

func ListTeachingLoadToProto(teachingLoad []domain.TeachingLoad) *journalv1.ListTeachingLoadResponse {
	return &journalv1.ListTeachingLoadResponse{
		TotalCount:    uint32(len(teachingLoad)),
		TeachingLoads: pkgmapper.MapSliceRef(teachingLoad, TeachingLoadToProto),
	}
}

func SubjectsToProto(subjects []domain.Subject) *journalv1.ListSubjectsResponse {

	res := pkgmapper.MapSliceRef(
		subjects,
		func(subject *domain.Subject) *journalv1.Subject {
			return &journalv1.Subject{
				Id:       subject.ID.String(),
				FullName: subject.FullName,
			}
		},
	)
	return &journalv1.ListSubjectsResponse{
		TotalCount: uint32(len(res)),
		Subjects:   res,
	}
}

// Support functions

func uuidPtrToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}
