package database

import "context"

type Transaction func(c context.Context, fn func(c context.Context) (err error)) error
