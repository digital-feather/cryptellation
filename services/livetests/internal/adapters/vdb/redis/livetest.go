package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	config "github.com/digital-feather/cryptellation/internal/go/adapters/redis"
	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

const (
	redisKeyLivetestIDs   = "livetests"
	redisKeyLivetest      = "livetest-%d"
	redisKeyMutexLivetest = "livetest-lock-%d"
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
		return nil, fmt.Errorf("loading redis config: %w", err)
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

func (db *DB) CreateLivetest(ctx context.Context, bt *livetest.Livetest) error {
	incr, err := db.client.Incr(ctx, redisKeyLivetestIDs).Result()
	if err != nil {
		return err
	}

	bt.ID = uint(incr)
	return db.client.Set(ctx, livetestKey(uint(incr)), bt, 0).Err()
}

func (db *DB) ReadLivetest(ctx context.Context, id uint) (livetest.Livetest, error) {
	bValue, err := db.client.Get(ctx, livetestKey(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			err = vdb.ErrRecordNotFound
		}
		return livetest.Livetest{}, err
	}

	bt := livetest.Livetest{}
	if err := json.Unmarshal(bValue, &bt); err != nil {
		return livetest.Livetest{}, err
	}

	return bt, nil
}

func (db *DB) UpdateLivetest(ctx context.Context, bt livetest.Livetest) error {
	return db.client.Set(ctx, livetestKey(bt.ID), bt, 0).Err()
}

func (db *DB) DeleteLivetest(ctx context.Context, bt livetest.Livetest) error {
	return db.client.Del(ctx, livetestKey(bt.ID)).Err()
}

func (db *DB) LockedLivetest(id uint, fn vdb.LockedLivetestCallback) error {
	mutex := db.lockClient.NewMutex(livetestMutexName(id), lockOptions...)
	if err := mutex.Lock(); err != nil {
		return err
	}

	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}

		ok, localErr := mutex.Unlock()
		if localErr != nil {
			err = localErr
		} else if !ok {
			err = fmt.Errorf("unlock failed for livetest %d", id)
		}
	}()

	err = fn()
	return err
}

func livetestKey(id uint) string {
	return fmt.Sprintf(redisKeyLivetest, id)
}

func livetestMutexName(id uint) string {
	return fmt.Sprintf(redisKeyMutexLivetest, id)
}
