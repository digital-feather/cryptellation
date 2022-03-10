package binance

import (
	"strconv"
	"time"

	binance "github.com/adshao/go-binance/v2"
	"github.com/cryptellation/cryptellation/pkg/types/candlestick"
	"github.com/cryptellation/cryptellation/pkg/types/period"
)

// TimeToKLineTime will take the time from a candle and will convert it to Kline time
func TimeToKLineTime(t time.Time) int64 {
	return t.UnixMilli()
}

// KLineTimeToTime will take the time from a kline and will convert it to candle time
func KLineTimeToTime(t int64) time.Time {
	return time.Unix(t/1000, 0)
}

// KLineToCandlestick will convert KLine binance format for Candlestick
func KLineToCandlestick(k binance.Kline, p period.Symbol, now time.Time) (time.Time, candlestick.Candlestick, error) {
	var c candlestick.Candlestick

	// Get time
	t := KLineTimeToTime(k.OpenTime)

	// Convert Open
	open, err := strconv.ParseFloat(k.Open, 64)
	if err != nil {
		return time.Unix(0, 0), c, err
	}

	// Convert High
	high, err := strconv.ParseFloat(k.High, 64)
	if err != nil {
		return time.Unix(0, 0), c, err
	}

	// Convert Low
	low, err := strconv.ParseFloat(k.Low, 64)
	if err != nil {
		return time.Unix(0, 0), c, err
	}

	// Convert Close
	close, err := strconv.ParseFloat(k.Close, 64)
	if err != nil {
		return time.Unix(0, 0), c, err
	}

	// Convert Volume
	volume, err := strconv.ParseFloat(k.Volume, 64)
	if err != nil {
		return time.Unix(0, 0), c, err
	}

	// Check completness
	uncomplete := false
	if now.Before(t.Add(p.Duration())) {
		uncomplete = true
	}

	// Instanciate Candle
	c = candlestick.Candlestick{
		Open:       open,
		High:       high,
		Low:        low,
		Close:      close,
		Volume:     volume,
		Uncomplete: uncomplete,
	}

	return t, c, nil
}

// KLinesToCandlesticks will transform a slice of binance format for Candlestick
func KLinesToCandlesticks(pair string, period period.Symbol, kl []*binance.Kline, now time.Time) (*candlestick.List, error) {
	cs := candlestick.NewList(candlestick.ListID{
		ExchangeName: Name,
		PairSymbol:   pair,
		Period:       period,
	})
	for _, k := range kl {
		t, c, err := KLineToCandlestick(*k, period, now)
		if err != nil {
			return nil, err
		}

		cs.Set(t, c)
		// TODO: Log warning in case of failure
	}

	return cs, nil
}
