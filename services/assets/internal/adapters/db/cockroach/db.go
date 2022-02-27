package cockroach

import (
	"context"

	"github.com/cryptellation/cryptellation/internal/adapters/cockroachdb"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
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
		return nil, err, func() {}
	}

	client, err := gorm.Open(postgres.Open(c.URL()), cockroachdb.DefaultGormConfig)
	if err != nil {
		return nil, err, func() {}
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

func (db *DB) CreateAssets(ctx context.Context, assets ...asset.Asset) error {
	entities := make([]Asset, len(assets))
	for i, model := range assets {
		entities[i].FromModel(model)
	}
	return db.client.WithContext(ctx).Create(&entities).Error
}

func (db *DB) ReadAssets(ctx context.Context, symbols ...string) ([]asset.Asset, error) {
	var ent []Asset
	if err := db.client.WithContext(ctx).Find(&ent, symbols).Error; err != nil {
		return nil, err
	}

	if len(ent) == 0 && len(symbols) != 0 {
		return nil, gorm.ErrRecordNotFound
	}

	models := make([]asset.Asset, len(ent))
	for i, entity := range ent {
		models[i] = entity.ToModel()
	}

	return models, nil
}

func (db *DB) UpdateAssets(ctx context.Context, assets ...asset.Asset) error {
	var entity Asset
	for _, model := range assets {
		entity.FromModel(model)

		if err := db.client.WithContext(ctx).Save(&entity).Error; err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) DeleteAssets(ctx context.Context, assets ...asset.Asset) error {
	var entity Asset
	for _, model := range assets {
		entity.FromModel(model)

		if err := db.client.WithContext(ctx).Delete(&entity).Error; err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Close() {
	sqlDb, err := db.client.DB()
	if err != nil {
		return
	}

	sqlDb.Close()
}

func Reset() error {
	db, err, closeDb := New()
	if err != nil {
		return err
	}
	defer closeDb()

	entities := []interface{}{
		&Asset{},
	}

	for _, entity := range entities {
		if err := db.client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error; err != nil {
			return err
		}
	}

	return nil
}
