package mapper

import (
	"errors"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errorToCode = map[string]codes.Code{
		"common.invalid_data":       codes.InvalidArgument,
		"common.internal_error":     codes.Internal,
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
