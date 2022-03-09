package binance

import (
	"context"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/cryptellation/cryptellation/internal/adapters/binance"
	"github.com/cryptellation/cryptellation/pkg/types/exchange"
	"github.com/cryptellation/cryptellation/pkg/types/pair"
	"github.com/cryptellation/cryptellation/services/exchanges/internal/adapters/exchanges"
	"golang.org/x/xerrors"
)

type ExchangeService struct {
	config binance.Config
	client *client.Client
}

func New() (*ExchangeService, error) {
	var c binance.Config
	if err := c.Load().Validate(); err != nil {
		return nil, xerrors.Errorf("loading cockroachdb config: %w", err)
	}

	return &ExchangeService{
		client: client.NewClient(
			c.ApiKey,
			c.SecretKey),
	}, nil
}

func (ps *ExchangeService) Infos(ctx context.Context) (exchange.Exchange, error) {
	exchangeInfos, err := ps.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return exchange.Exchange{}, err
	}

	pairSymbols := make([]string, len(exchangeInfos.Symbols))
	for i, bs := range exchangeInfos.Symbols {
		pairSymbols[i] = pair.Pair{
			BaseAssetSymbol:  bs.BaseAsset,
			QuoteAssetSymbol: bs.QuoteAsset,
		}.Symbol()
	}

	exch := exchanges.Binance
	exch.Pairs = pairSymbols
	exch.LastSyncTime = time.Now()

	return exch, nil
}
