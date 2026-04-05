package convert

import (
	"testing"
	"time"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateStudentToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *journalv1.CreateStudentRequest
		want    *domain.CreateStudent
		wantNil []string
	}{
		{
			name: "with all fields",
			input: &journalv1.CreateStudentRequest{
				ClassId:  toPtr("class-123"),
				FullName: "John Doe",
			},
			want: &domain.CreateStudent{
				ClassID:  idPtr("class-123"),
				FullName: "John Doe",
			},
		},
		{
			name: "with nil class_id",
			input: &journalv1.CreateStudentRequest{
				FullName: "John Doe",
			},
			want: &domain.CreateStudent{
				FullName: "John Doe",
			},
			wantNil: []string{"ClassID"},
		},
		{
			name:    "empty request",
			input:   &journalv1.CreateStudentRequest{},
			want:    &domain.CreateStudent{},
			wantNil: []string{"ClassID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := CreateStudentToDomain(tt.input)
			assert.Equal(t, tt.want.FullName, got.FullName)
			if tt.want.ClassID != nil {
				require.NotNil(t, got.ClassID)
				assert.Equal(t, *tt.want.ClassID, *got.ClassID)
			} else {
				assert.Nil(t, got.ClassID)
			}
		})
	}
}

func TestListStudentsToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      *journalv1.ListStudentsRequest
		wantClass  domain.ID // ← value type, не указатель
		wantLimit  uint32
		wantOffset uint32
	}{
		{
			name: "with all fields",
			input: &journalv1.ListStudentsRequest{
				ClassId: "class-123",
				Limit:   10,
				Offset:  20,
			},
			wantClass:  "class-123", // ← просто значение
			wantLimit:  10,
			wantOffset: 20,
		},
		{
			name:       "zero values",
			input:      &journalv1.ListStudentsRequest{},
			wantClass:  "", // ← пустое значение, не nil
			wantLimit:  0,
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ListStudentsToDomain(tt.input)

			assert.Equal(t, tt.wantClass, got.ClassID)
			assert.Equal(t, tt.wantLimit, got.Limit)
			assert.Equal(t, tt.wantOffset, got.Offset)
		})
	}
}

func TestUpdateStudentToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *journalv1.UpdateStudentRequest
		wantID  domain.ID
		wantNil []string
	}{
		{
			name: "with all fields",
			input: &journalv1.UpdateStudentRequest{
				Id:       "student-123",
				ClassId:  toPtr("class-123"),
				FullName: toPtr("John Doe"),
			},
			wantID: "student-123",
		},
		{
			name: "with nil optional fields",
			input: &journalv1.UpdateStudentRequest{
				Id: "student-123",
			},
			wantID:  "student-123",
			wantNil: []string{"ClassID", "FullName"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := UpdateStudentToDomain(tt.input)
			assert.Equal(t, tt.wantID, got.ID)
			if tt.input.ClassId != nil {
				require.NotNil(t, got.ClassID)
				assert.Equal(t, domain.ID(*tt.input.ClassId), *got.ClassID)
			} else {
				assert.Nil(t, got.ClassID)
			}
			assert.Equal(t, tt.input.FullName, got.FullName)
		})
	}
}

func TestListTeachersToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.ListTeachersRequest{
		Limit:  10,
		Offset: 20,
	}
	got := ListTeachersToDomain(req)

	assert.Equal(t, uint32(10), got.Limit)
	assert.Equal(t, uint32(20), got.Offset)
}

func TestUpdateTeacherToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.UpdateTeacherRequest{
		Id:       "teacher-123",
		FullName: "Jane Smith",
	}
	got := UpdateTeacherToDomain(req)

	assert.Equal(t, domain.ID("teacher-123"), got.ID)
	assert.Equal(t, "Jane Smith", got.FullName)
}

func TestCreateClassToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *journalv1.CreateClassRequest
		wantNil bool
	}{
		{
			name: "with graduation year",
			input: &journalv1.CreateClassRequest{
				ClassName:      "9A",
				YearOfStudy:    4,
				GraduationYear: toPtr(uint32(2030)),
			},
			wantNil: false,
		},
		{
			name: "without graduation year",
			input: &journalv1.CreateClassRequest{
				ClassName:   "9A",
				YearOfStudy: 4,
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := CreateClassToDomain(tt.input)
			assert.Equal(t, tt.input.ClassName, got.ClassName)
			assert.Equal(t, tt.input.YearOfStudy, got.YearOfStudy)
			if tt.wantNil {
				assert.Nil(t, got.GraduationYear)
			} else {
				assert.NotNil(t, got.GraduationYear)
				assert.Equal(t, tt.input.GraduationYear, got.GraduationYear)
			}
		})
	}
}

func TestListTeacherClassesToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.ListTeacherClassesRequest{
		TeacherId: "teacher-123",
		Limit:     10,
		Offset:    20,
	}
	got := ListTeacherClassesToDomain(req)

	assert.Equal(t, domain.ID("teacher-123"), got.TeacherID)
	assert.Equal(t, uint32(10), got.Limit)
	assert.Equal(t, uint32(20), got.Offset)
}

func TestUpdateClassToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *journalv1.UpdateClassRequest
	}{
		{
			name: "with all optional fields",
			input: &journalv1.UpdateClassRequest{
				Id:             "class-123",
				ClassName:      toPtr("9B"),
				YearOfStudy:    toPtr(uint32(3)),
				GraduationYear: toPtr(uint32(2029)),
			},
		},
		{
			name: "with nil optional fields",
			input: &journalv1.UpdateClassRequest{
				Id: "class-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := UpdateClassToDomain(tt.input)
			assert.Equal(t, domain.ID(tt.input.Id), got.ID)
			assert.Equal(t, tt.input.ClassName, got.ClassName)
			assert.Equal(t, tt.input.YearOfStudy, got.YearOfStudy)
			assert.Equal(t, tt.input.GraduationYear, got.GraduationYear)
		})
	}
}

func TestUpdateSubjectToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.UpdateSubjectRequest{
		Id:       "subject-123",
		FullName: "Mathematics",
	}
	got := UpdateSubjectToDomain(req)

	assert.Equal(t, domain.ID("subject-123"), got.ID)
	assert.Equal(t, "Mathematics", got.FullName)
}

// pkg/utils/convert/to_domain_test.go

func TestRecordGradeToDomain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	ts := timestamppb.New(now)

	tests := []struct {
		name               string
		input              *journalv1.RecordGradeRequest
		wantGrade          *uint32
		wantStatusCode     *domain.ID // ← *domain.ID, не *string!
		wantNilDateOfGrade bool
	}{
		{
			name: "with grade value and timestamp",
			input: &journalv1.RecordGradeRequest{
				SubjectId:    "subject-123",
				StudentId:    "student-123",
				DateOfGrade:  ts,
				Filter:       &journalv1.RecordGradeRequest_Grade{Grade: 5},
				LessonNumber: 3,
				Note:         toPtr("Good work"),
			},
			wantGrade: toPtr(uint32(5)),
		},
		{
			name: "with status code and timestamp",
			input: &journalv1.RecordGradeRequest{
				SubjectId:    "subject-123",
				StudentId:    "student-123",
				DateOfGrade:  ts,
				Filter:       &journalv1.RecordGradeRequest_StatusCodeId{StatusCodeId: "status-123"},
				LessonNumber: 3,
			},
			// ✅ Ожидаем *domain.ID, а не *string
			wantStatusCode: idPtr("status-123"),
		},
		// ... остальные кейсы
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := RecordGradeToDomain(tt.input)

			assert.Equal(t, domain.ID(tt.input.SubjectId), got.SubjectID)
			assert.Equal(t, domain.ID(tt.input.StudentId), got.StudentID)
			assert.Equal(t, tt.input.LessonNumber, got.LessonNumber)
			assert.Equal(t, tt.input.Note, got.Note)

			// ✅ Проверка oneof: Grade
			if tt.wantGrade != nil {
				require.NotNil(t, got.Grade)
				assert.Equal(t, *tt.wantGrade, *got.Grade)
				assert.Nil(t, got.StatusCodeID)
			}

			// ✅ Проверка oneof: StatusCodeID (*domain.ID)
			if tt.wantStatusCode != nil {
				require.NotNil(t, got.StatusCodeID)
				assert.Equal(t, *tt.wantStatusCode, *got.StatusCodeID) // domain.ID == domain.ID
				assert.Nil(t, got.Grade)
			}

			// ✅ Проверка optional timestamp
			if tt.wantNilDateOfGrade {
				assert.Nil(t, got.DateOfGrade)
			} else if tt.input.DateOfGrade != nil {
				assert.NotNil(t, got.DateOfGrade)
				assert.Equal(t, tt.input.DateOfGrade.AsTime().Unix(), got.DateOfGrade.Unix())
			}
		})
	}
}

func TestListGradesToDomain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	startTs := timestamppb.New(now)
	endTs := timestamppb.New(now.Add(24 * time.Hour))

	tests := []struct {
		name    string
		input   *journalv1.ListGradesRequest
		wantNil []string
		panics  bool
	}{
		{
			name: "with student_id filter",
			input: &journalv1.ListGradesRequest{
				SubjectId: toPtr("subject-123"),
				Filter:    &journalv1.ListGradesRequest_StudentId{StudentId: "student-123"},
				Start:     startTs,
				End:       endTs,
			},
		},
		{
			name: "with class_id filter",
			input: &journalv1.ListGradesRequest{
				Filter: &journalv1.ListGradesRequest_ClassId{ClassId: "class-123"},
			},
		},
		{
			name: "nil timestamps cause zero-time pointers",
			input: &journalv1.ListGradesRequest{
				Filter: &journalv1.ListGradesRequest_StudentId{StudentId: "student-123"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				assert.Panics(t, func() {
					ListGradesToDomain(tt.input)
				})
				return
			}
			got := ListGradesToDomain(tt.input)
			if tt.input.SubjectId != nil {
				assert.NotNil(t, got.SubjectID)
				assert.Equal(t, domain.ID(*tt.input.SubjectId), *got.SubjectID)
			}
			if v, ok := tt.input.Filter.(*journalv1.ListGradesRequest_StudentId); ok {
				assert.Equal(t, domain.ID(v.StudentId), got.StudentID)
			}
			if v, ok := tt.input.Filter.(*journalv1.ListGradesRequest_ClassId); ok {
				assert.Equal(t, domain.ID(v.ClassId), got.ClassID)
			}
			assert.NotNil(t, got.Start)
			assert.NotNil(t, got.End)
		})
	}
}

// pkg/utils/convert/to_domain_test.go
// pkg/utils/convert/to_domain_test.go
// pkg/utils/convert/to_domain_test.go
// pkg/utils/convert/to_domain_test.go

func TestUpdateGradeToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          *journalv1.UpdateGradeRequest
		wantGrade      *uint32
		wantStatusCode *domain.ID
		wantNote       *string
	}{
		{
			// ✅ Новый тест-кейс для покрытия ветки "Grade"
			name: "with grade value",
			input: &journalv1.UpdateGradeRequest{
				GradeId: "grade-123",
				Filter:  &journalv1.UpdateGradeRequest_Grade{Grade: 5}, // ← Grade variant
			},
			wantGrade: toPtr(uint32(5)), // ✅ ожидаем *uint32(5)
		},
		{
			name: "with status code and note",
			input: &journalv1.UpdateGradeRequest{
				GradeId: "grade-123",
				Filter:  &journalv1.UpdateGradeRequest_StatusCodeId{StatusCodeId: "absent"},
				Note:    toPtr("Updated comment"),
			},
			wantStatusCode: idPtr("absent"),
			wantNote:       toPtr("Updated comment"),
		},
		// ... остальные кейсы
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := UpdateGradeToDomain(tt.input)

			assert.Equal(t, domain.ID(tt.input.GradeId), got.GradeID)

			// ✅ Проверка oneof: Grade
			if tt.wantGrade != nil {
				require.NotNil(t, got.Grade)
				assert.Equal(t, *tt.wantGrade, *got.Grade)
				assert.Nil(t, got.StatusCodeID) // ← важно: должно быть nil
			}

			// ✅ Проверка oneof: StatusCodeID
			if tt.wantStatusCode != nil {
				require.NotNil(t, got.StatusCodeID)
				assert.Equal(t, *tt.wantStatusCode, *got.StatusCodeID)
				assert.Nil(t, got.Grade) // ← важно: должно быть nil
			}

			// ✅ Проверка optional Note
			assert.Equal(t, tt.input.Note, got.Note)
		})
	}
}

func TestRecordHomeworkToDomain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	startTs := timestamppb.New(now)
	endTs := timestamppb.New(now.Add(7 * 24 * time.Hour))

	tests := []struct {
		name         string
		input        *journalv1.RecordHomeworkRequest
		wantNilStart bool
		wantNilEnd   bool
	}{
		{
			name: "with all fields",
			input: &journalv1.RecordHomeworkRequest{
				TeacherId:       "teacher-123",
				ClassId:         "class-123",
				SubjectId:       "subject-123",
				DescriptionTask: "Solve problems 1-10",
				Start:           startTs,
				End:             endTs,
			},
		},
		{
			name: "nil timestamps produce nil pointers in domain",
			input: &journalv1.RecordHomeworkRequest{
				TeacherId:       "teacher-123",
				ClassId:         "class-123",
				SubjectId:       "subject-123",
				DescriptionTask: "Solve problems 1-10",
			},
			wantNilStart: true,
			wantNilEnd:   true,
		},
		{
			name: "only start is nil",
			input: &journalv1.RecordHomeworkRequest{
				TeacherId:       "teacher-123",
				ClassId:         "class-123",
				SubjectId:       "subject-123",
				DescriptionTask: "Solve problems 1-10",
				End:             endTs,
			},
			wantNilStart: true,
		},
		{
			name: "only end is nil",
			input: &journalv1.RecordHomeworkRequest{
				TeacherId:       "teacher-123",
				ClassId:         "class-123",
				SubjectId:       "subject-123",
				DescriptionTask: "Solve problems 1-10",
				Start:           startTs,
			},
			wantNilEnd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := RecordHomeworkToDomain(tt.input)

			assert.Equal(t, domain.ID(tt.input.TeacherId), got.TeacherID)
			assert.Equal(t, domain.ID(tt.input.ClassId), got.ClassID)
			assert.Equal(t, domain.ID(tt.input.SubjectId), got.SubjectID)
			assert.Equal(t, tt.input.DescriptionTask, got.DescriptionTask)

			if tt.wantNilStart {
				assert.Nil(t, got.Start, "Start should be nil when not set in proto")
			} else if tt.input.Start != nil {
				require.NotNil(t, got.Start)
				assert.Equal(t, tt.input.Start.AsTime().Unix(), got.Start.Unix())
			}

			if tt.wantNilEnd {
				assert.Nil(t, got.End, "End should be nil when not set in proto")
			} else if tt.input.End != nil {
				require.NotNil(t, got.End)
				assert.Equal(t, tt.input.End.AsTime().Unix(), got.End.Unix())
			}
		})
	}
}

func TestUpdateHomeworkToDomain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	startTs := timestamppb.New(now)
	endTs := timestamppb.New(now.Add(7 * 24 * time.Hour))

	tests := []struct {
		name         string
		input        *journalv1.UpdateHomeworkRequest
		wantNilStart bool
		wantNilEnd   bool
	}{
		{
			name: "with all fields",
			input: &journalv1.UpdateHomeworkRequest{
				Id:              "hw-123",
				DescriptionTask: "Updated task",
				Start:           startTs,
				End:             endTs,
			},
		},
		{
			name: "nil timestamps produce nil pointers in domain",
			input: &journalv1.UpdateHomeworkRequest{
				Id:              "hw-123",
				DescriptionTask: "Updated task",
			},
			wantNilStart: true,
			wantNilEnd:   true,
		},
		{
			name: "only start is nil",
			input: &journalv1.UpdateHomeworkRequest{
				Id:              "hw-123",
				DescriptionTask: "Updated task",
				End:             endTs,
			},
			wantNilStart: true,
		},
		{
			name: "only end is nil",
			input: &journalv1.UpdateHomeworkRequest{
				Id:              "hw-123",
				DescriptionTask: "Updated task",
				Start:           startTs,
			},
			wantNilEnd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := UpdateHomeworkToDomain(tt.input)

			assert.Equal(t, domain.ID(tt.input.Id), got.ID)
			assert.Equal(t, tt.input.DescriptionTask, got.DescriptionTask)

			if tt.wantNilStart {
				assert.Nil(t, got.Start, "Start should be nil when not set in proto")
			} else {
				require.NotNil(t, got.Start)
				assert.Equal(t, tt.input.Start.AsTime().Unix(), got.Start.Unix())
			}

			if tt.wantNilEnd {
				assert.Nil(t, got.End, "End should be nil when not set in proto")
			} else {
				require.NotNil(t, got.End)
				assert.Equal(t, tt.input.End.AsTime().Unix(), got.End.Unix())
			}
		})
	}
}

func TestListHomeworkToDomain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	startTs := timestamppb.New(now)
	endTs := timestamppb.New(now.Add(7 * 24 * time.Hour))

	tests := []struct {
		name         string
		input        *journalv1.ListHomeworkRequest
		wantNilStart bool
		wantNilEnd   bool
	}{
		{
			name: "with all fields",
			input: &journalv1.ListHomeworkRequest{
				ClassId:   "class-123",
				SubjectId: "subject-123",
				Start:     startTs,
				End:       endTs,
			},
		},
		{
			name: "nil timestamps produce nil pointers in domain",
			input: &journalv1.ListHomeworkRequest{
				ClassId:   "class-123",
				SubjectId: "subject-123",
			},
			wantNilStart: true,
			wantNilEnd:   true,
		},
		{
			name: "only start is nil",
			input: &journalv1.ListHomeworkRequest{
				ClassId:   "class-123",
				SubjectId: "subject-123",
				End:       endTs,
			},
			wantNilStart: true,
		},
		{
			name: "only end is nil",
			input: &journalv1.ListHomeworkRequest{
				ClassId:   "class-123",
				SubjectId: "subject-123",
				Start:     startTs,
			},
			wantNilEnd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ListHomeworkToDomain(tt.input)

			assert.Equal(t, domain.ID(tt.input.ClassId), got.ClassID)
			assert.Equal(t, domain.ID(tt.input.SubjectId), got.SubjectID)

			if tt.wantNilStart {
				assert.Nil(t, got.Start, "Start should be nil when not set in proto")
			} else {
				require.NotNil(t, got.Start)
				assert.Equal(t, tt.input.Start.AsTime().Unix(), got.Start.Unix())
			}

			if tt.wantNilEnd {
				assert.Nil(t, got.End, "End should be nil when not set in proto")
			} else {
				require.NotNil(t, got.End)
				assert.Equal(t, tt.input.End.AsTime().Unix(), got.End.Unix())
			}
		})
	}
}

func TestUpdateStatusCodeToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.UpdateStatusCodeRequest{
		Id:       "status-123",
		FullName: "Absent",
	}
	got := UpdateStatusCodeToDomain(req)

	assert.Equal(t, domain.ID("status-123"), got.ID)
	assert.Equal(t, "Absent", got.FullName)
}

func TestCreateTeachingLoadToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.CreateTeachingLoadRequest{
		TeacherId: "teacher-123",
		SubjectId: "subject-123",
		ClassId:   "class-123",
	}
	got := CreateTeachingLoadToDomain(req)

	assert.Equal(t, domain.ID("teacher-123"), got.TeacherID)
	assert.Equal(t, domain.ID("subject-123"), got.SubjectID)
	assert.Equal(t, domain.ID("class-123"), got.ClassID)
}

func TestUpdateTeachingLoadToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.UpdateTeachingLoadRequest{
		Id:        "load-123",
		TeacherId: "teacher-456",
	}
	got := UpdateTeachingLoadToDomain(req)

	assert.Equal(t, domain.ID("load-123"), got.ID)
	assert.Equal(t, domain.ID("teacher-456"), got.TeacherID)
}

func TestListTeachingLoadToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         *journalv1.ListTeachingLoadRequest
		wantTeacherID *domain.ID
		wantClassID   *domain.ID
	}{
		{
			name: "with teacher_id filter",
			input: &journalv1.ListTeachingLoadRequest{
				Filter: &journalv1.ListTeachingLoadRequest_TeacherId{TeacherId: "teacher-123"},
			},
			wantTeacherID: idPtr("teacher-123"),
			wantClassID:   nil,
		},
		{
			name: "with class_id filter",
			input: &journalv1.ListTeachingLoadRequest{
				Filter: &journalv1.ListTeachingLoadRequest_ClassId{ClassId: "class-123"},
			},
			wantTeacherID: nil,
			wantClassID:   idPtr("class-123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListTeachingLoadToDomain(tt.input)
			if tt.wantTeacherID != nil {
				require.NotNil(t, got.TeacherID)
				assert.Equal(t, *tt.wantTeacherID, *got.TeacherID)
			} else {
				assert.Nil(t, got.TeacherID)
			}
			if tt.wantClassID != nil {
				require.NotNil(t, got.ClassID)
				assert.Equal(t, *tt.wantClassID, *got.ClassID)
			} else {
				assert.Nil(t, got.ClassID)
			}
		})
	}
}

func toPtr[T any](a T) *T { return &a }

func idPtr(s string) *domain.ID {
	id := domain.ID(s)
	return &id
}
