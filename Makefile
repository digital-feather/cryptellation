.PHONY: proto
proto:
	@./scripts/proto.sh assets
	@./scripts/proto.sh candlesticks
	@./scripts/proto.sh exchanges
	@./scripts/proto.sh pairs