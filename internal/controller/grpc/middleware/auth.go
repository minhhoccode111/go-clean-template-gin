package middleware

import (
	"context"

	pkgjwt "github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const claimsKey ctxKey = "jwtClaims"

// ClaimsFromContext extracts the JWT claims from the context.
func ClaimsFromContext(ctx context.Context) (*pkgjwt.ClaimsJWT, bool) {
	claims, ok := ctx.Value(claimsKey).(*pkgjwt.ClaimsJWT)

	return claims, ok
}

// AuthInterceptor returns a gRPC unary interceptor for JWT authentication.
func AuthInterceptor(secret string) grpc.UnaryServerInterceptor {
	skipAuthMethods := map[string]bool{
		"/grpc.v1.AuthService/Register": true,
		"/grpc.v1.AuthService/Login":    true,
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if skipAuthMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		claims, err := pkgjwt.ValidateToken(values[0], secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		ctx = context.WithValue(ctx, claimsKey, claims)

		return handler(ctx, req)
	}
}
