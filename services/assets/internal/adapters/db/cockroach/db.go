package cockroach

import (
	"context"

	"github.com/cryptellation/cryptellation/internal/adapters/cockroachdb"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"github.com/cryptellation/cryptellation/services/assets/internal/adapters/db"
	"golang.org/x/xerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
	config cockroachdb.Config
}

func New() (*DB, error, func()) {
	var c cockroachdb.Config
	if err := c.Load().Validate(); err != nil {
		return nil, xerrors.Errorf("loading cockroachdb config: %w", err), func() {}
	}

	client, err := gorm.Open(postgres.Open(c.URL()), cockroachdb.DefaultGormConfig)
	if err != nil {
		return nil, xerrors.Errorf("opening cockroachdb connection: %w", err), func() {}
	}

	db := &DB{
		client: client,
		config: c,
	}

	closeFunc := func() {
		db.Close()
	}

	return db, nil, closeFunc
}

func (cockroach *DB) CreateAssets(ctx context.Context, assets ...asset.Asset) error {
	entities := make([]Asset, len(assets))
	for i, model := range assets {
		entities[i].FromModel(model)
	}

	err := cockroach.client.WithContext(ctx).Create(&entities).Error
	if err != nil {
		return xerrors.Errorf("creating %+v: %w", assets, err)
	}

	return nil
}

func (cockroach *DB) ReadAssets(ctx context.Context, symbols ...string) ([]asset.Asset, error) {
	var ent []Asset
	if err := cockroach.client.WithContext(ctx).Find(&ent, symbols).Error; err != nil {
		return nil, xerrors.Errorf("reading %+v: %w", symbols, err)
	}

	models := make([]asset.Asset, len(ent))
	for i, entity := range ent {
		models[i] = entity.ToModel()
	}

	return models, nil
}

func (cockroach *DB) UpdateAssets(ctx context.Context, assets ...asset.Asset) error {
	var entity Asset
	for _, model := range assets {
		entity.FromModel(model)

		if err := cockroach.client.WithContext(ctx).Save(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return xerrors.Errorf("updating %+v: %w", assets, db.ErrNotFound)
			}

			return xerrors.Errorf("updating %+v: %w", assets, err)
		}
	}
	return nil
}

func (cockroach *DB) DeleteAssets(ctx context.Context, assets ...asset.Asset) error {
	var entity Asset
	for _, model := range assets {
		entity.FromModel(model)

		if err := cockroach.client.WithContext(ctx).Delete(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return xerrors.Errorf("deleting %+v: %w", assets, db.ErrNotFound)
			}

			return xerrors.Errorf("deleting %+v: %w", assets, err)
		}
	}
	return nil
}

func (cockroach *DB) Close() {
	sqlDb, err := cockroach.client.DB()
	if err != nil {
		return
	}

	sqlDb.Close()
}

func Reset() error {
	db, err, closeDb := New()
	if err != nil {
		return xerrors.Errorf("creating connection for reset: %w", err)
	}
	defer closeDb()

	entities := []interface{}{
		&Asset{},
	}

	for _, entity := range entities {
		if err := db.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return xerrors.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}
