.PHONY: proto
proto:
	@./scripts/proto.sh backtests
	@./scripts/proto.sh candlesticks
	@./scripts/proto.sh exchanges