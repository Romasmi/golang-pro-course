package usecases

import (
	"context"
)

type HelloUsecase struct{}

func (u *HelloUsecase) Do(_ context.Context, _ any) (any, error) {
	return "hello world", nil
}
