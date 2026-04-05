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

// ============================================================================
// StudentService Tests
// ============================================================================

func TestCreateStudentToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *journalv1.CreateStudentRequest
		want    *domain.CreateStudent
		wantNil []string // поля, которые должны быть nil
	}{
		{
			name: "with all fields",
			input: &journalv1.CreateStudentRequest{
				ClassId:  toPtr("class-123"),
				FullName: "John Doe",
			},
			want: &domain.CreateStudent{
				ClassID:  toPtr("class-123"),
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
			assert.Equal(t, tt.want.ClassID, got.ClassID)
		})
	}
}

func TestListStudentsToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *journalv1.ListStudentsRequest
		want  *domain.ListStudentsRequest
	}{
		{
			name: "with all fields",
			input: &journalv1.ListStudentsRequest{
				ClassId: "class-123",
				Limit:   10,
				Offset:  20,
			},
			want: &domain.ListStudentsRequest{
				ClassID: "class-123",
				Limit:   10,
				Offset:  20,
			},
		},
		{
			name:  "zero values",
			input: &journalv1.ListStudentsRequest{},
			want:  &domain.ListStudentsRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ListStudentsToDomain(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateStudentToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   *journalv1.UpdateStudentRequest
		wantID  string
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
				assert.Equal(t, tt.input.ClassId, got.ClassID)
			} else {
				assert.Nil(t, got.ClassID)
			}
			if tt.input.FullName != nil {
				assert.Equal(t, tt.input.FullName, got.FullName)
			} else {
				assert.Nil(t, got.FullName)
			}
		})
	}
}

// ============================================================================
// TeacherService Tests
// ============================================================================

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

	assert.Equal(t, "teacher-123", got.ID)
	assert.Equal(t, "Jane Smith", got.FullName)
}

// ============================================================================
// ClassesService Tests
// ============================================================================

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

	assert.Equal(t, "teacher-123", got.TeacherID)
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
			assert.Equal(t, tt.input.Id, got.ID)
			assert.Equal(t, tt.input.ClassName, got.ClassName)
			assert.Equal(t, tt.input.YearOfStudy, got.YearOfStudy)
			assert.Equal(t, tt.input.GraduationYear, got.GraduationYear)
		})
	}
}

// ============================================================================
// SubjectService Tests
// ============================================================================

func TestUpdateSubjectToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.UpdateSubjectRequest{
		Id:       "subject-123",
		FullName: "Mathematics",
	}
	got := UpdateSubjectToDomain(req)

	assert.Equal(t, "subject-123", got.ID)
	assert.Equal(t, "Mathematics", got.FullName)
}

// ============================================================================
// GradesService Tests
// ============================================================================

func TestRecordGradeToDomain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	ts := timestamppb.New(now)

	tests := []struct {
		name               string
		input              *journalv1.RecordGradeRequest
		wantGrade          *uint32
		wantStatusCode     *string
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
			wantStatusCode: toPtr("status-123"),
		},
		{
			name: "nil timestamp produces nil pointer in domain",
			input: &journalv1.RecordGradeRequest{
				SubjectId:    "subject-123",
				StudentId:    "student-123",
				Filter:       &journalv1.RecordGradeRequest_Grade{Grade: 4},
				LessonNumber: 1,
			},
			wantGrade:          toPtr(uint32(4)),
			wantNilDateOfGrade: true,
		},
		{
			name: "nil filter leaves both grade and status nil",
			input: &journalv1.RecordGradeRequest{
				SubjectId: "subject-123",
				StudentId: "student-123",
			},
			wantNilDateOfGrade: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := RecordGradeToDomain(tt.input)

			assert.Equal(t, tt.input.SubjectId, got.SubjectID)
			assert.Equal(t, tt.input.StudentId, got.StudentID)
			assert.Equal(t, tt.input.LessonNumber, got.LessonNumber)
			assert.Equal(t, tt.input.Note, got.Note)

			if tt.wantGrade != nil {
				assert.NotNil(t, got.Grade)
				assert.Equal(t, *tt.wantGrade, *got.Grade)
				assert.Nil(t, got.StatusCodeID)
			} else if tt.wantStatusCode != nil {
				assert.NotNil(t, got.StatusCodeID)
				assert.Equal(t, *tt.wantStatusCode, *got.StatusCodeID)
				assert.Nil(t, got.Grade)
			} else {
				assert.Nil(t, got.Grade)
				assert.Nil(t, got.StatusCodeID)
			}

			if tt.wantNilDateOfGrade {
				assert.Nil(t, got.DateOfGrade, "DateOfGrade should be nil when not set in proto")
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
			name: "BUG: nil timestamps cause zero-time pointers",
			input: &journalv1.ListGradesRequest{
				Filter: &journalv1.ListGradesRequest_StudentId{StudentId: "student-123"},
				// Start and End are nil
			},
			// Note: Original code creates pointers to zero time, not nil
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
				assert.Equal(t, *tt.input.SubjectId, *got.SubjectID)
			}
			if v, ok := tt.input.Filter.(*journalv1.ListGradesRequest_StudentId); ok {
				assert.Equal(t, v.StudentId, got.StudentID)
			}
			if v, ok := tt.input.Filter.(*journalv1.ListGradesRequest_ClassId); ok {
				assert.Equal(t, v.ClassId, got.ClassID)
			}
			// Note: Original code always returns non-nil pointers for Start/End
			// even when input is nil (creates pointer to zero time)
			assert.NotNil(t, got.Start)
			assert.NotNil(t, got.End)
		})
	}
}
func TestUpdateGradeToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          *journalv1.UpdateGradeRequest
		wantGrade      *uint32
		wantStatusCode *string
		wantNote       *string
	}{
		{
			name: "with grade value",
			input: &journalv1.UpdateGradeRequest{
				GradeId: "grade-123",
				Filter:  &journalv1.UpdateGradeRequest_Grade{Grade: 5},
			},
			wantGrade: toPtr(uint32(5)),
		},
		{
			name: "with status code and note",
			input: &journalv1.UpdateGradeRequest{
				GradeId: "grade-123",
				Filter:  &journalv1.UpdateGradeRequest_StatusCodeId{StatusCodeId: "absent"},
				Note:    toPtr("Updated comment"),
			},
			wantStatusCode: toPtr("absent"),
			wantNote:       toPtr("Updated comment"),
		},
		{
			name: "nil filter leaves both pointers nil",
			input: &journalv1.UpdateGradeRequest{
				GradeId: "grade-123",
				// Filter not set
			},
			// wantGrade и wantStatusCode по умолчанию nil
		},
		{
			name: "nil note remains nil",
			input: &journalv1.UpdateGradeRequest{
				GradeId: "grade-123",
				Filter:  &journalv1.UpdateGradeRequest_Grade{Grade: 4},
				// Note not set
			},
			wantGrade: toPtr(uint32(4)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := UpdateGradeToDomain(tt.input)

			assert.Equal(t, tt.input.GradeId, got.GradeID)

			// 🔒 Безопасная проверка oneof с require перед разыменованием
			if tt.wantGrade != nil {
				require.NotNil(t, got.Grade, "Grade pointer should not be nil")
				assert.Equal(t, *tt.wantGrade, *got.Grade)
				assert.Nil(t, got.StatusCodeID, "StatusCodeID should be nil when grade is set")
			} else if tt.wantStatusCode != nil {
				require.NotNil(t, got.StatusCodeID, "StatusCodeID pointer should not be nil")
				assert.Equal(t, *tt.wantStatusCode, *got.StatusCodeID)
				assert.Nil(t, got.Grade, "Grade should be nil when status code is set")
			} else {
				assert.Nil(t, got.Grade, "Grade should be nil when filter not set")
				assert.Nil(t, got.StatusCodeID, "StatusCodeID should be nil when filter not set")
			}

			// Проверка optional Note
			if tt.wantNote != nil {
				require.NotNil(t, got.Note)
				assert.Equal(t, *tt.wantNote, *got.Note)
			} else {
				assert.Equal(t, tt.input.Note, got.Note) // оба должны быть nil или равны
			}
		})
	}
}

// ============================================================================
// HomeworkService Tests
// ============================================================================
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

			assert.Equal(t, tt.input.TeacherId, got.TeacherID)
			assert.Equal(t, tt.input.ClassId, got.ClassID)
			assert.Equal(t, tt.input.SubjectId, got.SubjectID)
			assert.Equal(t, tt.input.DescriptionTask, got.DescriptionTask)

			if tt.wantNilStart {
				assert.Nil(t, got.Start, "Start should be nil when not set in proto")
			} else if tt.input.Start != nil {
				require.NotNil(t, got.Start) // require перед разыменованием!
				assert.Equal(t, tt.input.Start.AsTime().Unix(), got.Start.Unix())
			}

			if tt.wantNilEnd {
				assert.Nil(t, got.End, "End should be nil when not set in proto")
			} else if tt.input.End != nil {
				require.NotNil(t, got.End) // require перед разыменованием!
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

			assert.Equal(t, tt.input.Id, got.ID)
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

			assert.Equal(t, tt.input.ClassId, got.ClassID)
			assert.Equal(t, tt.input.SubjectId, got.SubjectID)

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

// ============================================================================
// StatusCodeService Tests
// ============================================================================

func TestUpdateStatusCodeToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.UpdateStatusCodeRequest{
		Id:       "status-123",
		FullName: "Absent",
	}
	got := UpdateStatusCodeToDomain(req)

	assert.Equal(t, "status-123", got.ID)
	assert.Equal(t, "Absent", got.FullName)
}

// ============================================================================
// TeachingLoadService Tests
// ============================================================================

func TestCreateTeachingLoadToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.CreateTeachingLoadRequest{
		TeacherId: "teacher-123",
		SubjectId: "subject-123",
		ClassId:   "class-123",
	}
	got := CreateTeachingLoadToDomain(req)

	assert.Equal(t, "teacher-123", got.TeacherID)
	assert.Equal(t, "subject-123", got.SubjectID)
	assert.Equal(t, "class-123", got.ClassID)
}

func TestUpdateTeachingLoadToDomain(t *testing.T) {
	t.Parallel()

	req := &journalv1.UpdateTeachingLoadRequest{
		Id:        "load-123",
		TeacherId: "teacher-456",
	}
	got := UpdateTeachingLoadToDomain(req)

	assert.Equal(t, "load-123", got.ID)
	assert.Equal(t, "teacher-456", got.TeacherID)
}

func TestListTeachingLoadToDomain(t *testing.T) {
	t.Parallel()

	// ⚠️ BUG IN ORIGINAL CODE:
	// oneof handling always creates pointers to BOTH values,
	// even though only one should be set. This may break domain logic.

	tests := []struct {
		name          string
		input         *journalv1.ListTeachingLoadRequest
		wantTeacherID *string
		wantClassID   *string
	}{
		{
			name: "with teacher_id filter",
			input: &journalv1.ListTeachingLoadRequest{
				Filter: &journalv1.ListTeachingLoadRequest_TeacherId{TeacherId: "teacher-123"},
			},
			wantTeacherID: toPtr("teacher-123"),
			wantClassID:   toPtr(""), // Bug: returns pointer to empty string
		},
		{
			name: "with class_id filter",
			input: &journalv1.ListTeachingLoadRequest{
				Filter: &journalv1.ListTeachingLoadRequest_ClassId{ClassId: "class-123"},
			},
			wantTeacherID: toPtr(""), // Bug: returns pointer to empty string
			wantClassID:   toPtr("class-123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListTeachingLoadToDomain(tt.input)
			// Note: Due to bug, both pointers are always non-nil
			assert.NotNil(t, got.TeacherID)
			assert.NotNil(t, got.ClassID)
			assert.Equal(t, tt.wantTeacherID, got.TeacherID)
			assert.Equal(t, tt.wantClassID, got.ClassID)
		})
	}
}

func toPtr[T any](a T) *T { return &a }
