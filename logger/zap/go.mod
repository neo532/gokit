module github.com/neo532/gokit/logger/zap

go 1.26.1

require (
	github.com/neo532/gokit v0.0.0-00010101000000-000000000000
	github.com/neo532/gokit/logger/writer/lumberjack v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.28.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/neo532/gokit => ../..

replace github.com/neo532/gokit/logger/writer/lumberjack => ../writer/lumberjack
