package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptellation/cryptellation/pkg/types/exchange"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/adapters/db"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/adapters/exchanges"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/domain"
	"golang.org/x/xerrors"
)

type ReadAndSyncIfExpiredExchangesHandler struct {
	repository db.Port
	services   map[string]exchanges.Port
}

func NewReadAndSyncIfExpiredExchangesHandler(
	repository db.Port,
	services map[string]exchanges.Port,
) ReadAndSyncIfExpiredExchangesHandler {
	if repository == nil {
		panic("nil repository")
	}

	if services == nil || len(services) == 0 {
		panic("nil services")
	}

	return ReadAndSyncIfExpiredExchangesHandler{
		repository: repository,
		services:   services,
	}
}

func (reh ReadAndSyncIfExpiredExchangesHandler) Handle(
	ctx context.Context,
	expirationDuration *time.Duration,
	names ...string,
) ([]exchange.Exchange, error) {
	dbExchanges, err := reh.repository.ReadExchanges(ctx, names...)
	if err != nil {
		return nil, xerrors.Errorf("handling exchanges from db reading: %w", err)
	}

	toSync, err := domain.GetExpiredExchangesNames(names, dbExchanges, expirationDuration)
	if err != nil {
		return nil, xerrors.Errorf("determining exchanges to synchronize: %w", err)
	}

	synced, err := reh.getExchangeFromServices(ctx, toSync...)
	if err != nil {
		return nil, err
	}

	err = reh.upsertExchanges(ctx, dbExchanges, synced)
	if err != nil {
		return nil, err
	}

	mappedExchanges := exchange.ArrayToMap(dbExchanges)
	for _, exch := range synced {
		mappedExchanges[exch.Name] = exch
	}

	return exchange.MapToArray(mappedExchanges), nil
}

func (reh ReadAndSyncIfExpiredExchangesHandler) getExchangeFromServices(ctx context.Context, toSync ...string) ([]exchange.Exchange, error) {
	synced := make([]exchange.Exchange, 0, len(toSync))
	for _, name := range toSync {
		service, ok := reh.services[name]
		if !ok {
			return nil, xerrors.New(fmt.Sprintf("inexistant exchange service %q", name))
		}

		exch, err := service.Infos(ctx)
		if err != nil {
			return nil, err
		}

		synced = append(synced, exch)
	}

	return synced, nil
}

func (reh ReadAndSyncIfExpiredExchangesHandler) upsertExchanges(ctx context.Context, dbExchanges, toUpsert []exchange.Exchange) error {
	toCreate := make([]exchange.Exchange, 0, len(toUpsert))
	toUpdate := make([]exchange.Exchange, 0, len(toUpsert))
	mappedDbExchanges := exchange.ArrayToMap(dbExchanges)
	for _, exch := range toUpsert {
		if _, ok := mappedDbExchanges[exch.Name]; ok {
			toUpdate = append(toUpdate, exch)
		} else {
			toCreate = append(toCreate, exch)
		}
	}

	if len(toCreate) > 0 {
		if err := reh.repository.CreateExchanges(ctx, toCreate...); err != nil {
			return err
		}
	}

	if len(toUpdate) > 0 {
		if err := reh.repository.UpdateExchanges(ctx, toUpdate...); err != nil {
			return err
		}
	}

	return nil
}
