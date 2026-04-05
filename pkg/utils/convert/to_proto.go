package convert

import (
	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Студенты (Students)

func StudentToProto(params *domain.Student) *journalv1.Student {
	return &journalv1.Student{
		Id:       params.ID.String(),
		ClassId:  IDToProtoPtr(params.ClassID),
		FullName: params.FullName,
	}
}

func ListStudentsToProto(params *domain.ListStudentsResponse) *journalv1.ListStudentsResponse {
	return &journalv1.ListStudentsResponse{
		TotalCount: params.TotalCount,
		Students:   MapSlice(params.Students, StudentToProto),
	}
}

// Учителя (Teacher)

func TeacherToProto(params *domain.Teacher) *journalv1.Teacher {
	return &journalv1.Teacher{
		Id:       params.ID.String(),
		FullName: params.FullName,
	}
}

func ListTeachersToProto(params *domain.ListTeachersResponse) *journalv1.ListTeachersResponse {
	return &journalv1.ListTeachersResponse{
		TotalCount: params.TotalCount,
		Teachers:   MapSlice(params.Teachers, TeacherToProto),
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

func ListTeacherClassesToProto(params *domain.ListClasses) *journalv1.ListClassesResponse {
	return &journalv1.ListClassesResponse{
		TotalCount: params.TotalCount,
		Classes:    MapSlice(params.Classes, ClassToProto),
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
	} else if params.StatusCodeID != "" {
		res.Filter = &journalv1.Grade_StatusCodeId{
			StatusCodeId: params.StatusCodeID.String(),
		}
	}

	return res
}

func ListGradesToProto(params *domain.ListGradesResponse) *journalv1.ListGradesResponse {
	return &journalv1.ListGradesResponse{
		TotalCount: params.TotalCount,
		Grades:     MapSlice(params.Grades, GradeToProto),
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

func ListHomeworkToProto(params *domain.ListHomeworkResponse) *journalv1.ListHomeworkResponse {
	return &journalv1.ListHomeworkResponse{
		TotalCount: params.TotalCount,
		Homeworks:  MapSlice(params.Homeworks, HomeworkToProto),
	}
}

// Статус - коды (TeachingLoad)

func StatusCodeToProto(params *domain.StatusCode) *journalv1.StatusCode {
	return &journalv1.StatusCode{
		Id:       params.ID.String(),
		FullName: params.FullName,
	}
}

func ListStatusCodeToProto(params *domain.ListStatusCode) *journalv1.ListStatusCodeResponse {
	return &journalv1.ListStatusCodeResponse{
		TotalCount:  params.TotalCount,
		StatusCodes: MapSlice(params.StatusCodes, StatusCodeToProto),
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

func ListTeachingLoadToProto(params *domain.ListTeachingLoadResponse) *journalv1.ListTeachingLoadResponse {
	return &journalv1.ListTeachingLoadResponse{
		TotalCount:    params.TotalCount,
		TeachingLoads: MapSlice(params.TeachingLoads, TeachingLoadToProto),
	}
}

// Support functions

func MapSlice[F any, T any](items []F, mapper func(*F) T) []T {
	if items == nil {
		return nil
	}

	res := make([]T, len(items))

	for i, v := range items {
		res[i] = mapper(&v)
	}
	return res
}
