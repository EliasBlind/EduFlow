package auth

import (
	"context"

	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"github.com/EliasBlind/EduFlow/internal/sso_service/mapper"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Auth interface {
	Register(ctx context.Context, params *domain.RegisterRequest) error
	VerifyEmail(ctx context.Context, params *domain.VerifyRequest) (*uuid.UUID, error)
	Login(ctx context.Context, params *domain.LoginRequest) (*domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) (bool, error)
	RefreshToken(ctx context.Context, params *domain.RefreshRequest) (*domain.TokenPair, error)
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

func (s *serverAPI) VerifyEmail(ctx context.Context, req *ssov1.VerifyRequest) (*ssov1.VerifyResponse, error) {
	param := mapper.VerifyEmailToDomain(req)
	id, err := s.auth.VerifyEmail(ctx, param)
	if err != nil {
		return nil, err
	}
	return &ssov1.VerifyResponse{Id: id.String()}, nil
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
	token, err := s.auth.RefreshToken(ctx, param)
	if err != nil {
		return nil, err
	}
	return mapper.TokenPairToProto(token), nil
}
