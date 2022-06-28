.PHONY: all

proto:
	@./scripts/proto.sh backtests candlesticks exchanges ticks

lint:
	@./scripts/lint.sh
