package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	pkgmapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
	"github.com/EliasBlind/EduFlow/pkg/roles"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(s.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.ErrInternal
	}

	var id uuid.UUID
	if s.Id != nil && s.GetId() != "" {
		id, err = uuid.Parse(s.GetId())
		if err != nil {
			return nil, domain.ErrInvalidData
		}
	}

	role := roles.GetRole(s.GetRole()) // ← главное: парсим роль и берём адрес

	return &domain.User{
		Id:           id,
		Email:        s.GetEmail(),
		Login:        s.GetLogin(),
		PasswordHash: hash,
		Role:         &role,
	}, nil
}

func UsersToDomain(s *ssov1.CreateUsersRequest) ([]domain.User, error) {
	var firstErr error
	usrToDomain := func(u *ssov1.User) domain.User {
		usr, err := UserToDomain(u)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			return domain.User{} // не разыменовываем nil
		}
		return *usr
	}

	res := pkgmapper.MapSlice(s.Users, usrToDomain)
	if firstErr != nil {
		return nil, firstErr
	}
	return res, nil
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

func HashPassword(password string) ([]byte, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return hashedBytes, nil
}
