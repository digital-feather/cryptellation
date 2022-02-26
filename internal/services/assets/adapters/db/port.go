package db

import "github.com/cryptellation/cryptellation/pkg/types/asset"

type Port interface {
	Create(assets ...asset.Asset) error
	Read(symbols ...string) ([]asset.Asset, error)
	Update(assets ...asset.Asset) error
	Delete(assets ...asset.Asset) error
}
