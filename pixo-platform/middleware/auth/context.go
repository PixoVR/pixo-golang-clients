package auth

import (
	"context"
)

type CustomContext struct {
	FindUserByID func(id int) (*interface{}, error)
}

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(AbstractCustomContextKey).(*CustomContext)
	if !ok {
		return nil
	}

	return customContext
}
