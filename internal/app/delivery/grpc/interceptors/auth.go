package interceptors

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	AuthHeaderKey   = "Authorization"
	BearerTokenType = "Bearer"
)

func Auth(authenticator Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if isSecuredMethod(info.FullMethod) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "no metadata with auth token")
			}

			var authHeader string
			headers := md.Get(AuthHeaderKey)
			if len(headers) > 0 {
				authHeader = headers[0]
			}

			token, err := parseAuthHeader(authHeader)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}

			userID, err := authenticator.VerifyAccessToken(token)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}

			md.Append("user_id", strconv.FormatInt(userID, 10))
			ctx = metadata.NewIncomingContext(ctx, md)
		}
		return handler(ctx, req)
	}
}

func isSecuredMethod(name string) bool {
	return strings.HasSuffix(name, "UpdateUser") ||
		strings.HasSuffix(name, "GetUser") ||
		strings.HasSuffix(name, "ListUsers") ||
		strings.HasSuffix(name, "RemoveUser")
}

func parseAuthHeader(header string) (string, error) {
	authHeaderParts := strings.Split(header, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != BearerTokenType {
		return "", errors.New("invalid auth header")
	}
	return authHeaderParts[1], nil
}
