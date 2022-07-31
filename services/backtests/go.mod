module github.com/digital-feather/cryptellation/services/backtests

replace github.com/digital-feather/cryptellation/clients/go => ../../clients/go

replace github.com/digital-feather/cryptellation/services/candlesticks => ../candlesticks

replace github.com/digital-feather/cryptellation/internal/go => ../../internal/go

go 1.18

require (
	github.com/digital-feather/cryptellation/clients/go v0.0.0-00010101000000-000000000000
	github.com/digital-feather/cryptellation/internal/go v0.0.0-00010101000000-000000000000
	github.com/digital-feather/cryptellation/services/candlesticks v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-redsync/redsync/v4 v4.5.1
	github.com/stretchr/testify v1.8.0
	google.golang.org/grpc v1.48.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
