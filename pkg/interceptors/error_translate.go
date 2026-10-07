package interceptors

import (
	"context"
	"errors"
	"strings"

	"github.com/EliasBlind/EduFlow/pkg/i18n"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TranslatableError interface {
	error
	TranslationKey() string
}

func ErrorInterceptor(t *i18n.Translator, mapError func(error, *errdetails.LocalizedMessage) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		lang := extractLanguage(ctx)

		if transErr, ok := errors.AsType[TranslatableError](err); ok {
			translatedMsg := t.Translate(transErr.TranslationKey(), lang)

			details := &errdetails.LocalizedMessage{
				Locale:  lang,
				Message: translatedMsg,
			}

			return nil, mapError(err, details)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
}

func extractLanguage(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "en"
	}

	header := md.Get("accept-language")

	if len(header) > 0 {
		return strings.Split(header[0], ",")[0]
	}

	return "en"
}
