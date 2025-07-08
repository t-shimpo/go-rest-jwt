package auth

import (
	"context"
)

type contextKey string

const UserIDKey contextKey = "userID"

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (int, bool) {
	v := ctx.Value(UserIDKey)
	userID, ok := v.(int)
	return userID, ok
}
