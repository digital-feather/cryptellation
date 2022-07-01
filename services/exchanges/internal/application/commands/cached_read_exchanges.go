package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/db"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/domain/exchange"
	"golang.org/x/xerrors"
)

type CachedReadExchangesHandler struct {
	repository db.Port
	services   map[string]exchanges.Port
}

func NewCachedReadExchangesHandler(
	repository db.Port,
	services map[string]exchanges.Port,
) CachedReadExchangesHandler {
	if repository == nil {
		panic("nil repository")
	}

	if len(services) == 0 {
		panic("nil services")
	}

	return CachedReadExchangesHandler{
		repository: repository,
		services:   services,
	}
}

func (reh CachedReadExchangesHandler) Handle(
	ctx context.Context,
	expirationDuration *time.Duration,
	names ...string,
) ([]exchange.Exchange, error) {
	dbExchanges, err := reh.repository.ReadExchanges(ctx, names...)
	if err != nil {
		return nil, fmt.Errorf("handling exchanges from db reading: %w", err)
	}

	toSync, err := exchange.GetExpiredExchangesNames(names, dbExchanges, expirationDuration)
	if err != nil {
		return nil, fmt.Errorf("determining exchanges to synchronize: %w", err)
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

func (reh CachedReadExchangesHandler) getExchangeFromServices(ctx context.Context, toSync ...string) ([]exchange.Exchange, error) {
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

func (reh CachedReadExchangesHandler) upsertExchanges(ctx context.Context, dbExchanges, toUpsert []exchange.Exchange) error {
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
