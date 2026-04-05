package convert

import (
	"testing"
	"time"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStudentToProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  *domain.Student
		expect *journalv1.Student
	}{
		{
			name: "with all fields",
			input: &domain.Student{
				ID:       "stu-1",
				ClassID:  toPtr("class-1"),
				FullName: "John Doe",
			},
			expect: &journalv1.Student{
				Id:       "stu-1",
				ClassId:  toPtr("class-1"),
				FullName: "John Doe",
			},
		},
		{
			name: "with nil class_id",
			input: &domain.Student{
				ID:       "stu-2",
				FullName: "Jane Smith",
			},
			expect: &journalv1.Student{
				Id:       "stu-2",
				FullName: "Jane Smith",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := StudentToProto(tt.input)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestListStudentsToProto(t *testing.T) {
	t.Parallel()

	input := &domain.ListStudentsResponse{
		TotalCount: 2,
		Students: []domain.Student{
			{ID: "stu-1", FullName: "Alice"},
			{ID: "stu-2", FullName: "Bob"},
		},
	}
	got := ListStudentsToProto(input)

	assert.Equal(t, uint32(2), got.TotalCount)
	require.Len(t, got.Students, 2)
	assert.Equal(t, "stu-1", got.Students[0].Id)
	assert.Equal(t, "stu-2", got.Students[1].Id)
}

func TestTeacherToProto(t *testing.T) {
	t.Parallel()

	input := &domain.Teacher{
		ID:       "tea-1",
		FullName: "Prof. Smith",
	}
	got := TeacherToProto(input)

	assert.Equal(t, "tea-1", got.Id)
	assert.Equal(t, "Prof. Smith", got.FullName)
}

func TestListTeachersToProto(t *testing.T) {
	t.Parallel()

	input := &domain.ListTeachersResponse{
		TotalCount: 1,
		Teachers:   []domain.Teacher{{ID: "tea-1", FullName: "Prof. X"}},
	}
	got := ListTeachersToProto(input)

	assert.Equal(t, uint32(1), got.TotalCount)
	require.Len(t, got.Teachers, 1)
	assert.Equal(t, "tea-1", got.Teachers[0].Id)
}

func TestClassToProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  *domain.Class
		expect *journalv1.Class
	}{
		{
			name: "with graduation year",
			input: &domain.Class{
				ID:             "cls-1",
				ClassName:      "9A",
				YearOfStudy:    4,
				GraduationYear: uint32(2030),
			},
			expect: &journalv1.Class{
				Id:             "cls-1",
				ClassName:      "9A",
				YearOfStudy:    4,
				GraduationYear: uint32(2030),
			},
		},
		{
			name: "with nil graduation year",
			input: &domain.Class{
				ID:          "cls-2",
				ClassName:   "10B",
				YearOfStudy: 2,
			},
			expect: &journalv1.Class{
				Id:          "cls-2",
				ClassName:   "10B",
				YearOfStudy: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ClassToProto(tt.input)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestListTeacherClassesToProto(t *testing.T) {
	t.Parallel()

	input := &domain.ListClasses{
		TotalCount: 1,
		Classes:    []domain.Class{{ID: "cls-1", ClassName: "11A", YearOfStudy: 1}},
	}
	got := ListTeacherClassesToProto(input)

	assert.Equal(t, uint32(1), got.TotalCount)
	require.Len(t, got.Classes, 1)
	assert.Equal(t, "cls-1", got.Classes[0].Id)
}

func TestSubjectToProto(t *testing.T) {
	t.Parallel()

	input := &domain.Subject{
		ID:       "sub-1",
		FullName: "Mathematics",
	}
	got := SubjectToProto(input)

	assert.Equal(t, "sub-1", got.Id)
	assert.Equal(t, "Mathematics", got.FullName)
}

func TestGradeToProto(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name             string
		input            *domain.Grade
		wantGrade        *uint32
		wantStatusCodeID string
		wantNilFilter    bool
	}{
		{
			name: "with grade value",
			input: &domain.Grade{
				ID:           "gr-1",
				SubjectID:    "sub-1",
				StudentID:    "stu-1",
				ClassID:      "cls-1",
				DateOfGrade:  &now,
				Grade:        toPtr(uint32(5)),
				LessonNumber: toPtr(uint32(3)),
				Note:         toPtr("Excellent"),
			},
			wantGrade: toPtr(uint32(5)),
		},
		{
			name: "with status code",
			input: &domain.Grade{
				ID:           "gr-2",
				SubjectID:    "sub-1",
				StudentID:    "stu-1",
				ClassID:      "cls-1",
				StatusCodeID: "absent",
				LessonNumber: toPtr(uint32(1)),
			},
			wantStatusCodeID: "absent",
		},
		{
			name: "with nil optional fields",
			input: &domain.Grade{
				ID:        "gr-3",
				SubjectID: "sub-1",
				StudentID: "stu-1",
				ClassID:   "cls-1",
			},
			wantNilFilter: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GradeToProto(tt.input)

			assert.Equal(t, tt.input.ID, got.Id)
			assert.Equal(t, tt.input.SubjectID, got.SubjectId)
			assert.Equal(t, tt.input.StudentID, got.StudentId)
			assert.Equal(t, tt.input.ClassID, got.ClassId)
			assert.Equal(t, tt.input.Note, got.Note)

			if tt.input.LessonNumber != nil {
				assert.Equal(t, *tt.input.LessonNumber, got.LessonNumber)
			}

			if tt.input.DateOfGrade != nil {
				require.NotNil(t, got.DateOfGrade)
				assert.Equal(t, tt.input.DateOfGrade.Unix(), got.DateOfGrade.AsTime().Unix())
			}

			// ✅ Правильная проверка oneof через type switch
			if tt.wantGrade != nil {
				switch v := got.Filter.(type) {
				case *journalv1.Grade_Grade:
					assert.Equal(t, *tt.wantGrade, v.Grade)
				default:
					t.Errorf("expected Grade variant, got %T", got.Filter)
				}
			} else if tt.wantStatusCodeID != "" {
				switch v := got.Filter.(type) {
				case *journalv1.Grade_StatusCodeId:
					assert.Equal(t, tt.wantStatusCodeID, v.StatusCodeId)
				default:
					t.Errorf("expected StatusCodeId variant, got %T", got.Filter)
				}
			} else if tt.wantNilFilter {
				assert.Nil(t, got.Filter, "Filter should be nil when neither grade nor status code is set")
			}
		})
	}
}

func TestListGradesToProto(t *testing.T) {
	t.Parallel()

	input := &domain.ListGradesResponse{
		TotalCount: 1,
		Grades: []domain.Grade{
			{ID: "gr-1", SubjectID: "sub-1", StudentID: "stu-1", ClassID: "cls-1", Grade: toPtr(uint32((4)))},
		},
	}
	got := ListGradesToProto(input)

	assert.Equal(t, uint32(1), got.TotalCount)
	require.Len(t, got.Grades, 1)
	assert.Equal(t, "gr-1", got.Grades[0].Id)
}

func TestHomeworkToProto(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name    string
		input   *domain.Homework
		panics  bool
		wantNil []string
	}{
		{
			name: "with all fields",
			input: &domain.Homework{
				ID:              "hw-1",
				ClassID:         "cls-1",
				TeacherID:       "tea-1",
				SubjectID:       "sub-1",
				DescriptionTask: "Solve problems",
				Start:           &now,
				End:             &now,
			},
		},
		{
			name: "BUG: nil Start causes panic",
			input: &domain.Homework{
				ID:              "hw-2",
				ClassID:         "cls-1",
				TeacherID:       "tea-1",
				SubjectID:       "sub-1",
				DescriptionTask: "Solve problems",
				End:             &now,
			},
			panics: true,
		},
		{
			name: "BUG: nil End causes panic",
			input: &domain.Homework{
				ID:              "hw-3",
				ClassID:         "cls-1",
				TeacherID:       "tea-1",
				SubjectID:       "sub-1",
				DescriptionTask: "Solve problems",
				Start:           &now,
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				assert.Panics(t, func() {
					HomeworkToProto(tt.input)
				})
				return
			}
			got := HomeworkToProto(tt.input)
			assert.Equal(t, tt.input.ID, got.Id)
			assert.Equal(t, tt.input.ClassID, got.ClassId)
			assert.Equal(t, tt.input.TeacherID, got.TeacherId)
			assert.Equal(t, tt.input.SubjectID, got.SubjectId)
			assert.Equal(t, tt.input.DescriptionTask, got.DescriptionTask)
			if tt.input.Start != nil {
				assert.Equal(t, tt.input.Start.Unix(), got.Start.AsTime().Unix())
			}
			if tt.input.End != nil {
				assert.Equal(t, tt.input.End.Unix(), got.End.AsTime().Unix())
			}
		})
	}
}

func TestListHomeworkToProto(t *testing.T) {
	t.Parallel()

	now := time.Now()
	input := &domain.ListHomeworkResponse{
		TotalCount: 1,
		Homeworks: []domain.Homework{
			{
				ID: "hw-1", ClassID: "cls-1", TeacherID: "tea-1",
				SubjectID: "sub-1", DescriptionTask: "Task",
				Start: &now, End: &now,
			},
		},
	}
	got := ListHomeworkToProto(input)

	assert.Equal(t, uint32(1), got.TotalCount)
	require.Len(t, got.Homeworks, 1)
	assert.Equal(t, "hw-1", got.Homeworks[0].Id)
}

func TestStatusCodeToProto(t *testing.T) {
	t.Parallel()

	input := &domain.StatusCode{
		ID:       "st-1",
		FullName: "Absent",
	}
	got := StatusCodeToProto(input)

	assert.Equal(t, "st-1", got.Id)
	assert.Equal(t, "Absent", got.FullName)
}

func TestListStatusCodeToProto(t *testing.T) {
	t.Parallel()

	input := &domain.ListStatusCode{
		TotalCount: 2,
		StatusCodes: []domain.StatusCode{
			{ID: "st-1", FullName: "Present"},
			{ID: "st-2", FullName: "Absent"},
		},
	}
	got := ListStatusCodeToProto(input)

	assert.Equal(t, uint32(2), got.TotalCount)
	require.Len(t, got.StatusCodes, 2)
	assert.Equal(t, "st-1", got.StatusCodes[0].Id)
}

func TestTeachingLoadToProto(t *testing.T) {
	t.Parallel()

	input := &domain.TeachingLoad{
		ID:        "load-1",
		TeacherID: "tea-1",
		SubjectID: "sub-1",
		ClassID:   "cls-1",
	}
	got := TeachingLoadToProto(input)

	assert.Equal(t, "load-1", got.Id)
	assert.Equal(t, "tea-1", got.TeacherId)
	assert.Equal(t, "sub-1", got.SubjectId)
	assert.Equal(t, "cls-1", got.ClassId)
}

func TestListTeachingLoadToProto(t *testing.T) {
	t.Parallel()

	input := &domain.ListTeachingLoadResponse{
		TotalCount: 1,
		TeachingLoads: []domain.TeachingLoad{
			{ID: "load-1", TeacherID: "tea-1", SubjectID: "sub-1", ClassID: "cls-1"},
		},
	}
	got := ListTeachingLoadToProto(input)

	assert.Equal(t, uint32(1), got.TotalCount)
	require.Len(t, got.TeachingLoads, 1)
	assert.Equal(t, "load-1", got.TeachingLoads[0].Id)
}

func TestMapSlice(t *testing.T) {
	t.Parallel()

	t.Run("nil input returns nil", func(t *testing.T) {
		t.Parallel()
		var input []domain.Student
		got := MapSlice(input, StudentToProto)
		assert.Nil(t, got)
	})

	t.Run("empty slice returns empty slice", func(t *testing.T) {
		t.Parallel()
		input := []domain.Student{}
		got := MapSlice(input, StudentToProto)
		assert.Empty(t, got)
	})

	t.Run("maps single element", func(t *testing.T) {
		t.Parallel()
		input := []domain.Student{
			{ID: "stu-1", FullName: "Alice"},
		}
		got := MapSlice(input, StudentToProto)
		require.Len(t, got, 1)
		assert.Equal(t, "stu-1", got[0].Id)
		assert.Equal(t, "Alice", got[0].FullName)
	})

	t.Run("maps multiple elements", func(t *testing.T) {
		t.Parallel()
		input := []domain.Student{
			{ID: "stu-1", FullName: "Alice"},
			{ID: "stu-2", FullName: "Bob"},
			{ID: "stu-3", FullName: "Charlie"},
		}
		got := MapSlice(input, StudentToProto)
		require.Len(t, got, 3)
		assert.Equal(t, "stu-1", got[0].Id)
		assert.Equal(t, "stu-2", got[1].Id)
		assert.Equal(t, "stu-3", got[2].Id)
	})
}
