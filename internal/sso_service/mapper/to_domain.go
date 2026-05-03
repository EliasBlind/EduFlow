package mapper

import (
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
)

func RegisterToDomain(s *ssov1.RegisterRequest) *domain.RegisterRequest {
	return &domain.RegisterRequest{
		Email:    s.GetEmail(),
		Login:    s.GetLogin(),
		Password: s.GetPassword(),
		AppId:    int(s.GetAppId()),
	}
}

func LoginToDomain(s *ssov1.LoginRequest) *domain.LoginRequest {
	return &domain.LoginRequest{
		Login:    s.GetLogin(),
		Password: s.GetPassword(),
		AppId:    int(s.AppId),
	}
}

func VerifyEmailToDomain(s *ssov1.VerifyRequest) *domain.VerifyRequest {
	return &domain.VerifyRequest{
		Email: s.GetEmail(),
		Code:  s.GetCode(),
	}
}

func RefreshToDomain(s *ssov1.RefreshRequest) *domain.RefreshRequest {
	return &domain.RefreshRequest{
		RefreshToken: s.RefreshToken,
		AppId:        int(s.AppId),
	}
}
