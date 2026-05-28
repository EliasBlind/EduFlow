package mapper

import (
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/google/uuid"
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

func UserToDomain(s *ssov1.User) (*domain.User, error) {

	parsedUUID, err := uuid.Parse(s.Id)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	return &domain.User{
		Id:    parsedUUID,
		Email: s.Email,
		Login: s.Login,
	}, nil
}

func SetRoleToDomain(s *ssov1.SetRoleRequest) (*domain.User, error) {

	parsedUUID, err := uuid.Parse(s.UserId)
	if err != nil {
		return nil, domain.ErrInvalidData
	}

	role := roles.GetRole(s.Role)
	return &domain.User{
		Id:   parsedUUID,
		Role: &role,
	}, nil
}
