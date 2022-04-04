package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	config "github.com/digital-feather/cryptellation/internal/adapters/redis"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"golang.org/x/xerrors"
)

const (
	redisKeyBacktestIDs   = "backtests"
	redisKeyBacktest      = "backtest-%d"
	redisKeyMutexBacktest = "backtest-lock-%d"
)

var (
	lockOptions = []redsync.Option{
		redsync.WithExpiry(vdb.Expiration),
		redsync.WithRetryDelay(vdb.RetryDelay),
		redsync.WithTries(vdb.Tries),
	}
)

type DB struct {
	client     *redis.Client
	lockClient *redsync.Redsync
}

func New() (*DB, error) {
	var c config.Config
	if err := c.Load().Validate(); err != nil {
		return nil, xerrors.Errorf("loading redis config: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.Password, // no password set
		DB:       0,          // use default DB
	})

	// TODO Check connection

	pool := goredis.NewPool(client)
	ls := redsync.New(pool)

	return &DB{
		client:     client,
		lockClient: ls,
	}, nil
}

func (db *DB) CreateBacktest(ctx context.Context, bt *backtest.Backtest) error {
	incr, err := db.client.Incr(ctx, redisKeyBacktestIDs).Result()
	if err != nil {
		return err
	}

	bt.ID = uint(incr)
	return db.client.Set(ctx, backtestKey(uint(incr)), bt, 0).Err()
}

func (db *DB) ReadBacktest(ctx context.Context, id uint) (backtest.Backtest, error) {
	bValue, err := db.client.Get(ctx, backtestKey(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			err = vdb.ErrRecordNotFound
		}
		return backtest.Backtest{}, err
	}

	bt := backtest.Backtest{}
	if err := json.Unmarshal(bValue, &bt); err != nil {
		return backtest.Backtest{}, err
	}

	return bt, nil
}

func (db *DB) UpdateBacktest(ctx context.Context, bt backtest.Backtest) error {
	return db.client.Set(ctx, backtestKey(bt.ID), bt, 0).Err()
}

func (db *DB) DeleteBacktest(ctx context.Context, bt backtest.Backtest) error {
	return db.client.Del(ctx, backtestKey(bt.ID)).Err()
}

func (db *DB) LockedBacktest(id uint, fn vdb.LockedBacktestCallback) error {
	mutex := db.lockClient.NewMutex(backtestMutexName(id), lockOptions...)
	if err := mutex.Lock(); err != nil {
		return err
	}

	var err error
	defer func() {
		recover()
		ok, localErr := mutex.Unlock()
		if localErr != nil {
			err = localErr
		} else if !ok {
			err = errors.New("Unlock failed")
		}
	}()

	err = fn()
	return err
}

func backtestKey(id uint) string {
	return fmt.Sprintf(redisKeyBacktest, id)
}

func backtestMutexName(id uint) string {
	return fmt.Sprintf(redisKeyMutexBacktest, id)
}
