module github.com/neo532/gokit/database/orm

go 1.26.1

replace github.com/neo532/gokit => ../..

require (
	github.com/neo532/gokit v0.0.0-00010101000000-000000000000
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.20.0 // indirect
)
