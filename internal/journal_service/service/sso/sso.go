package sso

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	ssov1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/sso/v1"
)

type Auth struct {
	log     *slog.Logger
	cfg     *config.SsoConfig
	conn    *grpc.ClientConn
	ssoRole ssov1.AuthClient
}

func New(
	log *slog.Logger,
	cfg *config.SsoConfig,
) *Auth {
	host := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to grpc server", slog.Any("err", err))
		panic("Error connect")
	}

	return &Auth{
		log:     log,
		cfg:     cfg,
		conn:    conn,
		ssoRole: ssov1.NewAuthClient(conn),
	}
}

func (a *Auth) CreateUserSso(
	ctx context.Context,
	jwt string,
	user *domain.User,
) error {
	md := metadata.Pairs("authorization", "Bearer "+jwt)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := a.ssoRole.CreateUsers(
		ctx,
		&ssov1.CreateUsersRequest{
			Users: []*ssov1.User{
				{
					Id:    proto.String(user.ID.String()),
					Login: user.Login,
					Email: user.Email,
					Role:  user.Role.String(),
				},
			},
		},
	)

	if err != nil {
		a.log.Error("journalv1.CreateStudent grpc call failed", slog.Any("err", err))
		return domain.ErrInternal
	}

	return nil
}

func (a *Auth) Stop() {
	if err := a.conn.Close(); err != nil {
		a.log.Error("failed to close grpc connection", slog.Any("err", err))
	}
}
