package service

import (
	"context"
	"time"

	"github.com/cryptellation/cryptellation/internal/genproto/candlesticks"
	"google.golang.org/grpc"
)

type MockedCandlesticksClient struct {
}

func (m MockedCandlesticksClient) ReadCandlesticks(
	ctx context.Context,
	in *candlesticks.ReadCandlesticksRequest,
	opts ...grpc.CallOption,
) (*candlesticks.ReadCandlesticksResponse, error) {
	start, err := time.Parse(time.RFC3339, in.Start)
	if err != nil {
		return nil, err
	}

	return &candlesticks.ReadCandlesticksResponse{
		Candlesticks: []*candlesticks.Candlestick{
			{
				Time:   start.Format(time.RFC3339),
				Open:   1,
				High:   2,
				Low:    0.5,
				Close:  1.5,
				Volume: 500,
			},
		},
	}, nil
}
