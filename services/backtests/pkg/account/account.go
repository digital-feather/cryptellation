package account

import "errors"

var (
	ErrInvalidBalanceAmount = errors.New("invalid-balance-amount")
	ErrInvalidBalanceAsset  = errors.New("invalid-balance-asset")
)

type Account struct {
	Balances map[string]float64
}

func (a Account) Validate() error {
	for asset, balance := range a.Balances {
		if asset == "" {
			return ErrInvalidBalanceAsset
		}

		if balance < 0 {
			return ErrInvalidBalanceAmount
		}
	}

	return nil
}
