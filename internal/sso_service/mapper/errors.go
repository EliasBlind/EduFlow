package mapper

import (
	"errors"

	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
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

func MapToRPCError(err error, details *errdetails.LocalizedMessage) error {
	var appErr domain.AppError

	if !errors.As(err, &appErr) {
		return status.Error(codes.Internal, "internal error")
	}
	code := getGrpcCode(appErr.Key)

	st := status.New(code, details.Message)

	descSt, errDetail := st.WithDetails(details)
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
