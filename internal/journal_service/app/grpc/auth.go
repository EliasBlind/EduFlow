package grpcapp

import (
	"context"
	"log/slog"
	"strings"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TokenParser interface {
	ParseToken(token string) (*domain.UserClaims, error)
}

func UnaryAuthInterceptor(parser TokenParser, log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		const op = "interceptors.UnaryAuthInterceptor"

		log := log.With(
			slog.String("op", op),
			slog.String("method", info.FullMethod),
		)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Warn("unauthenticated: metadata is missing")
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			log.Warn("unauthenticated: authorization header is missing")
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		parts := strings.Split(authHeader[0], " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn("unauthenticated: invalid auth format")
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		claims, err := parser.ParseToken(parts[1])
		if err != nil {
			log.Warn("unauthenticated: token parsing failed", slog.String("error", err.Error()))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		log.Info("user authenticated",
			slog.String("user_id", claims.ID.String()),
			slog.String("role", claims.Role.String()),
		)

		newCtx := domain.ContextWithClaims(ctx, claims)

		return handler(newCtx, req)
	}
}
