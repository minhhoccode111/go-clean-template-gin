package middleware_test

import (
	"context"
	"testing"
	"time"

	grpcmw "github.com/minhhoccode111/go-clean-template-gin/internal/controller/grpc/middleware"
	pkgjwt "github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	respOK          = "ok"
	testSecret      = "test-secret"
	testUserID      = "42"
	grpcGetTaskPath = "/grpc.v1.TaskService/GetTask"
)

type ctxCapture struct {
	ctx context.Context
}

func (c *ctxCapture) handler(_ context.Context, _ any) (any, error) {
	return respOK, nil
}

func (c *ctxCapture) capturingHandler(ctx context.Context, _ any) (any, error) {
	c.ctx = ctx

	return respOK, nil
}

func genToken(t *testing.T) string {
	t.Helper()

	manager := pkgjwt.New(testSecret, time.Hour)

	token, err := manager.GenerateToken(testUserID)
	require.NoError(t, err)

	return token
}

func runSkipAuthTest(t *testing.T, method string) {
	t.Helper()

	jwtManager := pkgjwt.New(testSecret, time.Hour)
	interceptor := grpcmw.AuthInterceptor(jwtManager)
	info := &grpc.UnaryServerInfo{FullMethod: method}

	called := false
	handler := func(_ context.Context, _ any) (any, error) {
		called = true

		return respOK, nil
	}

	resp, err := interceptor(t.Context(), nil, info, handler)

	require.NoError(t, err)
	assert.Equal(t, respOK, resp)
	assert.True(t, called)
}

func TestAuthInterceptor_SkipRegister(t *testing.T) {
	t.Parallel()
	runSkipAuthTest(t, "/grpc.v1.AuthService/Register")
}

func TestAuthInterceptor_SkipLogin(t *testing.T) {
	t.Parallel()
	runSkipAuthTest(t, "/grpc.v1.AuthService/Login")
}

func TestAuthInterceptor_MissingMetadata(t *testing.T) {
	t.Parallel()

	jwtManager := pkgjwt.New(testSecret, time.Hour)
	interceptor := grpcmw.AuthInterceptor(jwtManager)
	info := &grpc.UnaryServerInfo{FullMethod: grpcGetTaskPath}

	capture := &ctxCapture{}

	resp, err := interceptor(t.Context(), nil, info, capture.handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "missing metadata")
}

func TestAuthInterceptor_MissingAuthorizationToken(t *testing.T) {
	t.Parallel()

	jwtManager := pkgjwt.New(testSecret, time.Hour)
	interceptor := grpcmw.AuthInterceptor(jwtManager)
	info := &grpc.UnaryServerInfo{FullMethod: grpcGetTaskPath}

	md := metadata.New(map[string]string{"other-key": "value"})
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	resp, err := interceptor(ctx, nil, info, capture.handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "missing authorization token")
}

func TestAuthInterceptor_InvalidToken(t *testing.T) {
	t.Parallel()

	jwtManager := pkgjwt.New(testSecret, time.Hour)
	interceptor := grpcmw.AuthInterceptor(jwtManager)
	info := &grpc.UnaryServerInfo{FullMethod: grpcGetTaskPath}

	md := metadata.Pairs("authorization", "Bearer invalid-token")
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	resp, err := interceptor(ctx, nil, info, capture.handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "invalid or expired token")
}

func TestAuthInterceptor_ValidToken(t *testing.T) {
	t.Parallel()

	jwtManager := pkgjwt.New(testSecret, time.Hour)
	interceptor := grpcmw.AuthInterceptor(jwtManager)
	info := &grpc.UnaryServerInfo{FullMethod: grpcGetTaskPath}

	token := genToken(t)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	resp, err := interceptor(ctx, nil, info, capture.capturingHandler)

	require.NoError(t, err)
	assert.Equal(t, respOK, resp)

	userID, ok := grpcmw.UserIDFromContext(capture.ctx)
	assert.True(t, ok)
	assert.Equal(t, testUserID, userID)
}

func TestUserIDFromContext_WithValue(t *testing.T) {
	t.Parallel()

	jwtManager := pkgjwt.New(testSecret, time.Hour)
	interceptor := grpcmw.AuthInterceptor(jwtManager)
	info := &grpc.UnaryServerInfo{FullMethod: grpcGetTaskPath}

	token := genToken(t)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	_, err := interceptor(ctx, nil, info, capture.capturingHandler)
	require.NoError(t, err)

	userID, ok := grpcmw.UserIDFromContext(capture.ctx)
	assert.True(t, ok)
	assert.Equal(t, testUserID, userID)
}

func TestUserIDFromContext_WithoutValue(t *testing.T) {
	t.Parallel()

	userID, ok := grpcmw.UserIDFromContext(t.Context())
	assert.False(t, ok)
	assert.Empty(t, userID)
}
