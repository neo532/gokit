package redis

/*
 * @abstract Redis pooler
 * @mail neo532@126.com
 * @date 2024-05-18
 */

import (
	"context"
	"math/rand"

	"github.com/go-redis/redis/v8"
)

type Pooler interface {
	Choose(c context.Context, rdbs *Rdbs) *redis.Client
}

type RandomPolicy struct {
}

func (*RandomPolicy) Choose(c context.Context, rdbs *Rdbs) *redis.Client {
	l := len(rdbs.rdbs)
	if l == 1 {
		return rdbs.rdbs[0].Client()
	}
	return rdbs.rdbs[rand.Intn(l)].Client()
}
