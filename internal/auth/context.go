package auth

import (
	"context"
	"errors"
)

type contextKey string

const userKey contextKey = "auth_user"

type UserClaims struct {
	Sub string
	Email string
	Role string
}

func WithUser(ctx context.Context, user *UserClaims) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (*UserClaims, error) {
	user, ok := ctx.Value(userKey).(*UserClaims)
	if !ok {
		return nil, errors.New("user not found in context.")
	}
	return user, nil
}