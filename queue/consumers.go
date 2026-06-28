package queue

import (
	"context"
	"strings"
)

type Consumers struct {
	csm map[string]Consumer
}

func NewConsumers(cs ...Consumer) Consumer {
	groups := &Consumers{
		csm: make(map[string]Consumer, len(cs)),
	}
	for _, o := range cs {
		groups.csm[o.Name()] = o
	}
	return groups
}
func (cs *Consumers) Start(ctx context.Context) (err error) {
	for _, o := range cs.csm {
		if e := o.Start(ctx); e != nil {
			err = e
		}
	}
	return
}
func (cs *Consumers) Stop(ctx context.Context) (err error) {
	for _, o := range cs.csm {
		if e := o.Stop(ctx); e != nil {
			err = e
		}
	}
	return
}
func (cs *Consumers) Name() (name string) {
	for _, o := range cs.csm {
		name += "," + o.Name()
	}
	return strings.TrimPrefix(name, ",")
}
