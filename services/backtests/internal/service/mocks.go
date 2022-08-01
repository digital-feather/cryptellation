package service

import (
	"context"
	"time"

	candlesticksProto "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
	"google.golang.org/grpc"
)

type MockedCandlesticksClient struct {
}

func (m MockedCandlesticksClient) ReadCandlesticks(
	ctx context.Context,
	in *candlesticksProto.ReadCandlesticksRequest,
	opts ...grpc.CallOption,
) (*candlesticksProto.ReadCandlesticksResponse, error) {
	start, err := time.Parse(time.RFC3339, in.Start)
	if err != nil {
		return nil, err
	}

	return &candlesticksProto.ReadCandlesticksResponse{
		Candlesticks: []*candlesticksProto.Candlestick{
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
