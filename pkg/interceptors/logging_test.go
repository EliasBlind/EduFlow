package interceptors

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type recordCaptureHandler struct {
	records []slog.Record
	mu      sync.Mutex
}

func (h *recordCaptureHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *recordCaptureHandler) Handle(ctx context.Context, rec slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = append(h.records, rec)
	return nil
}

func (h *recordCaptureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *recordCaptureHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *recordCaptureHandler) Records() []slog.Record {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]slog.Record{}, h.records...)
}

func TestUnaryServerInterceptor_Success(t *testing.T) {
	t.Parallel()

	handler := &recordCaptureHandler{}
	logger := slog.New(handler)
	interceptor := UnaryServerInterceptor(logger)

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	callHandler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	resp, err := interceptor(context.Background(), "req", info, callHandler)

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
	assert.Empty(t, handler.Records())
}

func TestUnaryServerInterceptor_Error(t *testing.T) {
	t.Parallel()

	handler := &recordCaptureHandler{}
	logger := slog.New(handler)
	interceptor := UnaryServerInterceptor(logger)

	mockErr := errors.New("test error")
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	callHandler := func(ctx context.Context, req any) (any, error) {
		return nil, mockErr
	}

	resp, err := interceptor(context.Background(), "req", info, callHandler)

	assert.Equal(t, mockErr, err)
	assert.Nil(t, resp)

	records := handler.Records()
	assert.Len(t, records, 1)
	assert.Equal(t, slog.LevelError, records[0].Level)

	attrs := make(map[string]any)
	records[0].Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})
	assert.Equal(t, "/test.Service/Method", attrs["method"])
	assert.Equal(t, mockErr, attrs["err"])
}

func TestUnaryServerInterceptor_ContextPropagation(t *testing.T) {
	t.Parallel()

	logger := slog.Default()
	interceptor := UnaryServerInterceptor(logger)

	type ctxKey string
	key := ctxKey("test")
	expected := "value"

	ctx := context.WithValue(context.Background(), key, expected)
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

	callHandler := func(ctx context.Context, req any) (any, error) {
		assert.Equal(t, expected, ctx.Value(key))
		return "ok", nil
	}

	resp, err := interceptor(ctx, "req", info, callHandler)

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestUnaryServerInterceptor_NilLogger(t *testing.T) {
	t.Parallel()

	interceptor := UnaryServerInterceptor(nil)
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	callHandler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	assert.NotPanics(t, func() {
		resp, err := interceptor(context.Background(), "req", info, callHandler)
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp)
	})
}
