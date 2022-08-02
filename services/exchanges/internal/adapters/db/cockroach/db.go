package cockroach

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/db"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
	config Config
}

func New() (*DB, error) {
	var c Config
	if err := c.Load().Validate(); err != nil {
		return nil, fmt.Errorf("loading cockroachdb config: %w", err)
	}

	client, err := gorm.Open(postgres.Open(c.URL()), DefaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening cockroachdb connection: %w", err)
	}

	db := &DB{
		client: client,
		config: c,
	}

	return db, nil
}

func (cockroach *DB) CreateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	entities := make([]Exchange, len(exchanges))
	for i, model := range exchanges {
		entities[i].FromModel(model)
	}

	err := cockroach.client.WithContext(ctx).Create(&entities).Error
	if err != nil {
		return fmt.Errorf("creating %+v: %w", exchanges, err)
	}

	return nil
}

func (cockroach *DB) ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	var ent []Exchange
	if err := cockroach.client.WithContext(ctx).Preload("Pairs").Preload("Periods").Find(&ent, names).Error; err != nil {
		return nil, fmt.Errorf("reading %+v: %w", names, err)
	}

	models := make([]exchange.Exchange, len(ent))
	for i, entity := range ent {
		models[i] = entity.ToModel()
	}

	return models, nil
}

func (cockroach *DB) UpdateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	var entity Exchange
	for _, model := range exchanges {
		entity.FromModel(model)

		if err := cockroach.client.WithContext(ctx).Updates(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("updating %+v: %w", exchanges, db.ErrNotFound)
			}

			return fmt.Errorf("updating %+v: %w", exchanges, err)
		}

		if err := cockroach.client.WithContext(ctx).Model(&entity).Association("Pairs").Replace(entity.Pairs); err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("replacing pairs associations from %+v: %w", exchanges, db.ErrNotFound)
			}

			return fmt.Errorf("replacing pairs associations from %+v: %w", exchanges, err)
		}

		if err := cockroach.client.WithContext(ctx).Model(&entity).Association("Periods").Replace(entity.Periods); err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("replacing periods associations from %+v: %w", exchanges, db.ErrNotFound)
			}

			return fmt.Errorf("replacing periods associations from %+v: %w", exchanges, err)
		}
	}
	return nil
}

func (cockroach *DB) DeleteExchanges(ctx context.Context, names ...string) error {
	for _, n := range names {
		if err := cockroach.client.WithContext(ctx).Delete(&Exchange{
			Name: n,
		}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("deleting %+v: %w", names, db.ErrNotFound)
			}

			return fmt.Errorf("deleting %+v: %w", names, err)
		}

		err := cockroach.client.WithContext(ctx).
			Where("NOT EXISTS(SELECT NULL FROM exchanges_pairs ep WHERE ep.pair_symbol = symbol)").
			Delete(&Pair{}).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("deleting unlinked pairs %+v: %w", names, db.ErrNotFound)
			}

			return fmt.Errorf("deleting unlinked pairs %+v: %w", names, err)
		}

		err = cockroach.client.WithContext(ctx).
			Where("NOT EXISTS(SELECT NULL FROM exchanges_periods ep WHERE ep.period_symbol = symbol)").
			Delete(&Period{}).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("deleting unlinked periods %+v: %w", names, db.ErrNotFound)
			}

			return fmt.Errorf("deleting unlinked periods %+v: %w", names, err)
		}
	}
	return nil
}

func (cockroach *DB) Reset() error {
	entities := []interface{}{
		&Exchange{},
		&Pair{},
		&Period{},
	}

	for _, entity := range entities {
		if err := cockroach.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return fmt.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
