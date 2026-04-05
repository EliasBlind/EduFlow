package classes

import (
	"context"
	"errors"
	"testing"

	"io"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ClassesMock struct {
	mock.Mock
}

func (m *ClassesMock) CreateClass(ctx context.Context, params *domain.CreateClass) (*domain.Class, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Class), args.Error(1)
}

func (m *ClassesMock) GetClass(ctx context.Context, class_id string) (*domain.Class, error) {
	args := m.Called(ctx, class_id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Class), args.Error(1)
}

func (m *ClassesMock) ListTeacherClasses(ctx context.Context, params *domain.ListTeacherClasses) (*domain.ListClasses, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListClasses), args.Error(1)
}

func (m *ClassesMock) UpdateClass(ctx context.Context, params *domain.UpdateClass) (*domain.Class, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Class), args.Error(1)
}

func (m *ClassesMock) DeleteClassRequest(ctx context.Context, class_id string) error {
	args := m.Called(ctx, class_id)
	return args.Error(0)
}

func TestServerAPI_AllMethods(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockClasses := new(ClassesMock)

	gRPCServer := grpc.NewServer()
	Register(gRPCServer, logger, mockClasses)

	api := &serverAPI{
		classes: mockClasses,
		log:     logger,
	}

	ctx := context.Background()
	testErr := errors.New("database error")

	t.Run("CreateClass", func(t *testing.T) {
		req := &journalv1.CreateClassRequest{}
		req.ClassName = "Math"

		mockClasses.On("CreateClass", ctx, mock.Anything).Return(&domain.Class{ID: "1"}, nil).Once()
		res, err := api.CreateClass(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		mockClasses.On("CreateClass", ctx, mock.Anything).Return(nil, testErr).Once()
		res, err = api.CreateClass(ctx, req)
		assert.Error(t, err)
	})

	t.Run("GetClass", func(t *testing.T) {
		req := &journalv1.GetClassRequest{Id: "1"}

		mockClasses.On("GetClass", ctx, "1").Return(&domain.Class{ID: "1"}, nil).Once()
		res, err := api.GetClass(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "1", res.Id)

		mockClasses.On("GetClass", ctx, "1").Return(nil, testErr).Once()
		_, err = api.GetClass(ctx, req)
		assert.Error(t, err)
	})

	t.Run("ListTeacherClasses", func(t *testing.T) {
		req := &journalv1.ListTeacherClassesRequest{TeacherId: "t1"}

		mockClasses.On("ListTeacherClasses", ctx, mock.Anything).Return(&domain.ListClasses{}, nil).Once()
		res, err := api.ListTeacherClasses(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		mockClasses.On("ListTeacherClasses", ctx, mock.Anything).Return(nil, testErr).Once()
		_, err = api.ListTeacherClasses(ctx, req)
		assert.Error(t, err)
	})

	t.Run("UpdateClass", func(t *testing.T) {
		req := &journalv1.UpdateClassRequest{}
		req.Id = "1"
		req.ClassName = proto.String("Physics")

		mockClasses.On("UpdateClass", ctx, mock.Anything).Return(&domain.Class{ID: "1"}, nil).Once()
		res, err := api.UpdateClass(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		mockClasses.On("UpdateClass", ctx, mock.Anything).Return(nil, testErr).Once()
		_, err = api.UpdateClass(ctx, req)
		assert.Error(t, err)
	})

	t.Run("DeleteClass", func(t *testing.T) {
		req := &journalv1.DeleteClassRequest{Id: "1"}

		mockClasses.On("DeleteClassRequest", ctx, "1").Return(nil).Once()
		res, err := api.DeleteClass(ctx, req)
		assert.NoError(t, err)
		assert.Nil(t, res)

		mockClasses.On("DeleteClassRequest", ctx, "1").Return(testErr).Once()
		_, err = api.DeleteClass(ctx, req)
		assert.Error(t, err)
	})
}
