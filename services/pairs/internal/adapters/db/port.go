package db

import (
	"context"

	"github.com/cryptellation/cryptellation/pkg/types/pair"
)

type Port interface {
	CreatePairs(ctx context.Context, pairs ...pair.Pair) error
	ReadPairs(ctx context.Context, pairs ...pair.Pair) ([]pair.Pair, error)
	UpdatePairs(ctx context.Context, pairs ...pair.Pair) error
	DeletePairs(ctx context.Context, pairs ...pair.Pair) error
}
