package auth

import (
	"context"

	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"github.com/EliasBlind/EduFlow/internal/sso_service/mapper"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Auth interface {
	Register(ctx context.Context, params *domain.RegisterRequest) error
	VerifyEmail(ctx context.Context, params *domain.VerifyRequest) (*domain.TokenPair, error)
	Login(ctx context.Context, params *domain.LoginRequest) (*domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) (bool, error)
	RefreshToken(ctx context.Context, params *domain.RefreshRequest) (*domain.TokenPair, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	SetRole(ctx context.Context, user *domain.User) error
	CreateStudent(ctx context.Context, user *domain.User) error
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(
		gRPC,
		&serverAPI{
			auth: auth,
		},
	)
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	param := mapper.RegisterToDomain(req)
	err := s.auth.Register(ctx, param)
	if err != nil {
		return nil, err
	}

	return &ssov1.RegisterResponse{}, nil
}

func (s *serverAPI) VerifyEmail(ctx context.Context, req *ssov1.VerifyRequest) (*ssov1.TokenPair, error) {
	param := mapper.VerifyEmailToDomain(req)
	token, err := s.auth.VerifyEmail(ctx, param)
	if err != nil {
		return nil, err
	}
	return mapper.TokenPairToProto(token), nil
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.TokenPair, error) {
	param := mapper.LoginToDomain(req)

	loginRequest, err := s.auth.Login(ctx, param)
	if err != nil {
		return nil, err
	}
	return mapper.TokenPairToProto(loginRequest), nil
}

func (s *serverAPI) Logout(ctx context.Context, req *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	success, err := s.auth.Logout(ctx, req.RefreshToken)
	if err != nil {
		return &ssov1.LogoutResponse{Success: false}, err
	}
	return &ssov1.LogoutResponse{Success: success}, nil
}

func (s *serverAPI) RefreshToken(ctx context.Context, req *ssov1.RefreshRequest) (*ssov1.TokenPair, error) {
	param := mapper.RefreshToDomain(req)
	refreshToken, err := s.auth.RefreshToken(ctx, param)
	if err != nil {
		return nil, err
	}
	return mapper.TokenPairToProto(refreshToken), nil
}

func (s *serverAPI) ListUsers(ctx context.Context, req *ssov1.ListUsersRequest) (*ssov1.ListUsersResponse, error) {
	users, err := s.auth.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.UsersToProto(users), nil
}

func (s *serverAPI) SetRole(ctx context.Context, req *ssov1.SetRoleRequest) (*emptypb.Empty, error) {
	user, err := mapper.SetRoleToDomain(req)
	if err != nil {
		return nil, err
	}
	err = s.auth.SetRole(ctx, user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *serverAPI) CreateStudent(ctx context.Context, req *ssov1.User) (*emptypb.Empty, error) {
	user, err := mapper.UserToDomain(req)
	if err != nil {
		return nil, err
	}

	err = s.auth.CreateStudent(ctx, user)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
