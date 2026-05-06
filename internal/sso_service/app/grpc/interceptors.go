package grpcapp

import (
	"context"

	"github.com/EliasBlind/EduFlow/internal/sso_service/mapper"
	"github.com/EliasBlind/EduFlow/pkg/i18n"
	"google.golang.org/grpc"
)

func ErrorInterceptor(t *i18n.Translator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)

		if err == nil {
			return resp, nil
		}

		return resp, mapper.MapToRPCError(ctx, err, t)
	}
}
