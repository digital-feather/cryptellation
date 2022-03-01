package cockroach

import (
	"context"

	"github.com/cryptellation/cryptellation/internal/adapters/cockroachdb"
	"github.com/cryptellation/cryptellation/pkg/types/pair"
	"github.com/cryptellation/cryptellation/services/pairs/internal/adapters/db"
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

func (cockroach *DB) CreatePairs(ctx context.Context, pairs ...pair.Pair) error {
	entities := make([]Pair, len(pairs))
	for i, model := range pairs {
		entities[i].FromModel(model)
	}

	err := cockroach.client.WithContext(ctx).Create(&entities).Error
	if err != nil {
		return xerrors.Errorf("creating %+v: %w", pairs, err)
	}

	return nil
}

func (cockroach *DB) ReadPairs(ctx context.Context, pairs ...pair.Pair) ([]pair.Pair, error) {
	var ent []Pair

	var tx *gorm.DB
	if len(pairs) == 0 {
		tx = cockroach.client
	} else {
		keys := make([][]interface{}, len(pairs))
		for i, p := range pairs {
			keys[i] = []interface{}{
				p.BaseSymbol,
				p.QuoteSymbol,
			}
		}

		tx = cockroach.client.Where("(base_symbol, quote_symbol) IN ?", keys)
	}

	if err := tx.Find(&ent).Error; err != nil {
		return nil, xerrors.Errorf("reading %+v: %w", pairs, err)
	}

	models := make([]pair.Pair, len(ent))
	for i, entity := range ent {
		models[i] = entity.ToModel()
	}

	return models, nil
}

func (cockroach *DB) UpdatePairs(ctx context.Context, pairs ...pair.Pair) error {
	var entity Pair
	for _, model := range pairs {
		entity.FromModel(model)

		if err := cockroach.client.WithContext(ctx).Save(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return xerrors.Errorf("updating %+v: %w", pairs, db.ErrNotFound)
			}

			return xerrors.Errorf("updating %+v: %w", pairs, err)
		}
	}
	return nil
}

func (cockroach *DB) DeletePairs(ctx context.Context, pairs ...pair.Pair) error {
	var entity Pair
	for _, model := range pairs {
		entity.FromModel(model)

		if err := cockroach.client.WithContext(ctx).Delete(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return xerrors.Errorf("deleting %+v: %w", pairs, db.ErrNotFound)
			}

			return xerrors.Errorf("deleting %+v: %w", pairs, err)
		}
	}
	return nil
}

func Reset() error {
	db, err := New()
	if err != nil {
		return xerrors.Errorf("creating connection for reset: %w", err)
	}

	entities := []interface{}{
		&Pair{},
	}

	for _, entity := range entities {
		if err := db.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return xerrors.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
