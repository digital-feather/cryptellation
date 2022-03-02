package db

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/pair"
)

type Port interface {
	CreatePairs(ctx context.Context, pairs ...pair.Pair) error
	ReadPairs(ctx context.Context, symbols ...string) ([]pair.Pair, error)
	UpdatePairs(ctx context.Context, pairs ...pair.Pair) error
	DeletePairs(ctx context.Context, symbols ...string) error
}
