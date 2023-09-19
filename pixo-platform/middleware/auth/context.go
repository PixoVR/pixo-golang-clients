package auth

import (
	"context"
)

type CustomContext struct {
	Services     *interface{}
	FindUserByID func(id int) (*interface{}, error)
}

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(CustomContextKey).(*CustomContext)
	if !ok {
		return nil
	}

	return customContext
}
