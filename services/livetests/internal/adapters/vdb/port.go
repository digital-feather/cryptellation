package vdb

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type LockedLivetestCallback func() error

type Port interface {
	CreateLivetest(ctx context.Context, bt *livetest.Livetest) error
	ReadLivetest(ctx context.Context, id uint) (livetest.Livetest, error)
	UpdateLivetest(ctx context.Context, bt livetest.Livetest) error
	DeleteLivetest(ctx context.Context, bt livetest.Livetest) error

	LockedLivetest(id uint, fn LockedLivetestCallback) error
}
