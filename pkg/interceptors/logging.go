package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			log.Error("grpc error",
				slog.String("method", info.FullMethod),
				slog.Any("err", err),
			)
		}
		return resp, err
	}
}
