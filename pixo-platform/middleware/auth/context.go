package auth

import (
	"context"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

type CustomContext struct {
	FindUserByID func(id int) (*platform.User, error)
}

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(CustomContextKey).(*CustomContext)
	if !ok {
		return nil
	}

	return customContext
}
