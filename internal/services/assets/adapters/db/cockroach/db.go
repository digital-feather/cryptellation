package cockroach

import (
	"github.com/cryptellation/cryptellation/internal/services/assets/adapters/db/cockroach/entities"
	"github.com/cryptellation/cryptellation/internal/utils/adapters/cockroachdb"
	"github.com/cryptellation/cryptellation/pkg/types/asset"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	client *gorm.DB
	config cockroachdb.Config
}

func New(c cockroachdb.Config) (*DB, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	client, err := gorm.Open(postgres.Open(c.URL()), cockroachdb.DefaultGormConfig)
	if err != nil {
		return nil, err
	}

	return &DB{
		client: client,
		config: c,
	}, nil
}

func (db *DB) Create(assets ...asset.Asset) error {
	entities := make([]entities.Asset, len(assets))
	for i, model := range assets {
		entities[i].FromModel(model)
	}
	return db.client.Create(&entities).Error
}

func (db *DB) Read(symbols ...string) ([]asset.Asset, error) {
	var ent []entities.Asset
	if err := db.client.Find(&ent, symbols).Error; err != nil {
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

func (db *DB) Update(assets ...asset.Asset) error {
	var entity entities.Asset
	for _, model := range assets {
		entity.FromModel(model)

		if err := db.client.Save(&entity).Error; err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Delete(assets ...asset.Asset) error {
	var entity entities.Asset
	for _, model := range assets {
		entity.FromModel(model)

		if err := db.client.Delete(&entity).Error; err != nil {
			return err
		}
	}
	return nil
}
