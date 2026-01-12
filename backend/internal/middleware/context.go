package middleware

import "context"

const UserIDKey contextKey = "user_id"

func GetUserID(ctx context.Context) int {
	if userID, ok := ctx.Value(UserIDKey).(int); ok {
		return userID
	}
	return 0
}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}
