package usecases

import (
	"context"
)

type Usecase interface {
	Do(ctx context.Context, req any) (any, error)
}
