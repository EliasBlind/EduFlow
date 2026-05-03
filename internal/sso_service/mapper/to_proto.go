package mapper

import (
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
)

func TokenPairToProto(s *domain.TokenPair) *ssov1.TokenPair {
	return &ssov1.TokenPair{
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
	}
}
