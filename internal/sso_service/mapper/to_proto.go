package mapper

import (
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	pkgmapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
)

func TokenPairToProto(s *domain.TokenPair) *ssov1.TokenPair {
	return &ssov1.TokenPair{
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
	}
}

func UsersToProto(users []domain.User) *ssov1.ListUsersResponse {
	result := pkgmapper.MapSliceRef(users, func(user *domain.User) *ssov1.User {
		id := user.Id.String()
		return &ssov1.User{
			Id:    &id,
			Login: user.Login,
			Email: user.Email,
			Role:  user.Role.String(),
		}
	})

	return &ssov1.ListUsersResponse{
		TotalCount: uint32(len(result)),
		Users:      result,
	}
}
