package mapper

import (
	"context"
	"errors"
	"strings"

	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"github.com/EliasBlind/EduFlow/pkg/i18n"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errorToCode = map[string]codes.Code{
		"common.invalid_data":       codes.InvalidArgument,
		"common.internal_error":     codes.Internal,
		"auth.weak_password":        codes.InvalidArgument,
		"auth.user_already_exists":  codes.AlreadyExists,
		"auth.user_not_found":       codes.NotFound,
		"auth.unauthenticated":      codes.Unauthenticated,
		"auth.access_denied":        codes.PermissionDenied,
		"verify.invalid_code":       codes.InvalidArgument,
		"verify.code_expired":       codes.DeadlineExceeded,
		"verify.code_not_found":     codes.NotFound,
		"infra.storage_unavailable": codes.Unavailable,
	}
)

func MapToRPCError(ctx context.Context, err error, trans *i18n.Translator) error {
	if err == nil {
		return nil
	}

	var appErr domain.AppError
	if !errors.As(err, &appErr) {
		return status.Error(codes.Internal, "internal error")
	}

	code := getGrpcCode(appErr.Key)

	lang := extractLanguage(ctx)

	translatedMsg := trans.Translate(appErr.Key, lang)

	st := status.New(code, translatedMsg)

	localizedDetail := &errdetails.LocalizedMessage{
		Locale:  lang,
		Message: translatedMsg,
	}

	descSt, errDetail := st.WithDetails(localizedDetail)
	if errDetail != nil {
		return st.Err()
	}

	return descSt.Err()
}

func getGrpcCode(key string) codes.Code {
	if code, ok := errorToCode[key]; ok {
		return code
	}
	return codes.Internal
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
