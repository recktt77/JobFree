package middleware

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryAuthInterceptor(allowedRoles ...string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, errors.New("missing token")
		}

		token := strings.TrimPrefix(tokens[0], "Bearer ")
		role, err := validateJWT(token)
		if err != nil {
			return nil, err
		}

		// Check role
		for _, r := range allowedRoles {
			if role == r {
				return handler(ctx, req)
			}
		}
		return nil, errors.New("access denied: role not allowed")
	}
}

// Fake validation
func validateJWT(token string) (string, error) {
	// return "client", nil // test
	if token == "admin-token" {
		return "admin", nil
	} else if token == "client-token" {
		return "client", nil
	} else if token == "freelancer-token" {
		return "freelancer", nil
	}
	return "", errors.New("invalid token")
}
