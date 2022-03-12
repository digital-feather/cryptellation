package cockroach

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/internal/adapters/cockroachdb"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/adapters/db"
	"github.com/cryptellation/cryptellation/services/candlesticks/internal/domain/candlestick"
	"golang.org/x/xerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
	config cockroachdb.Config
}

func New() (*DB, error) {
	var c cockroachdb.Config
	if err := c.Load().Validate(); err != nil {
		return nil, xerrors.Errorf("loading cockroachdb config: %w", err)
	}

	client, err := gorm.Open(postgres.Open(c.URL()), cockroachdb.DefaultGormConfig)
	if err != nil {
		return nil, xerrors.Errorf("opening cockroachdb connection: %w", err)
	}

	db := &DB{
		client: client,
		config: c,
	}

	return db, nil
}

func (d *DB) CreateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE := FromModelListToEntityList(cs)
	return d.client.WithContext(ctx).Create(&listCE).Error
}

func (d *DB) ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
	tx := d.client.Where(`
		exchange_name = ? AND
		pair_symbol = ? AND
		period_symbol = ? AND
		time BETWEEN ? AND ?`,
		cs.ExchangeName(),
		cs.PairSymbol(),
		cs.Period().String(),
		start, end)

	if limit != 0 {
		tx = tx.Limit(int(limit))
	}

	cse := []Candlestick{}
	if err := tx.WithContext(ctx).Find(&cse).Error; err != nil {
		return err
	}

	for _, ce := range cse {
		_, _, _, t, m := ce.ToModel()
		if err := cs.Set(t, m); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) UpdateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE := FromModelListToEntityList(cs)
	for _, ce := range listCE {
		tx := d.client.WithContext(ctx).
			Model(&Candlestick{}).
			Where("exchange_name = ?", ce.ExchangeName).
			Where("pair_symbol = ?", ce.PairSymbol).
			Where("period_symbol = ?", ce.PeriodSymbol).
			Where("time = ?", ce.Time).
			Updates(ce)

		if tx.Error != nil {
			return xerrors.Errorf("updating candlestick %q: %w", ce.Time, tx.Error)
		} else if tx.RowsAffected == 0 {
			return db.ErrNotFound
		}
	}

	return nil
}

func (d *DB) DeleteCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE := FromModelListToEntityList(cs)
	return d.client.WithContext(ctx).Delete(&listCE).Error
}

func Reset() error {
	db, err := New()
	if err != nil {
		return xerrors.Errorf("creating connection for reset: %w", err)
	}

	entities := []interface{}{
		&Candlestick{},
	}

	for _, entity := range entities {
		if err := db.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return xerrors.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
