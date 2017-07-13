package context

import (
	"context"

	"lenslocked.com/models"
)

const (
	userKey privateKey = "signed_in_user"
)

type privateKey string

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) (*models.User, bool) {
	u, ok := ctx.Value(userKey).(*models.User)
	return u, ok
}
