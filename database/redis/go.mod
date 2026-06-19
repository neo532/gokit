module github.com/neo532/gokit/database/redis

go 1.26.1

replace github.com/neo532/gokit => ../..

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/neo532/gokit v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
