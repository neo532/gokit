module github.com/neo532/gokit/database/redis

go 1.23.2

replace github.com/neo532/gokit => ../..

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/neo532/gokit v1.0.43
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
