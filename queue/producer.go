package queue

import (
	"context"

	"github.com/neo532/gokit/errorx"
)

type ProducerOption func(*Producers)

func WithProducer(pdc Producer) ProducerOption {
	return func(p *Producers) {
		p.def = pdc
	}
}
func WithProducerShadow(pdc Producer) ProducerOption {
	return func(p *Producers) {
		p.shadow = pdc
	}
}

func WithProducerGray(pdc Producer) ProducerOption {
	return func(p *Producers) {
		p.gray = pdc
	}
}

type Producers struct {
	def    Producer
	shadow Producer
	gray   Producer

	err error

	isGrayer Grayer
	isShadow Benchmarker
}

func NewProducers(ps ...ProducerOption) Producer {
	pdc := &Producers{
		isGrayer: &DefaultGrayer{},
		isShadow: &DefaultBenchmarker{},
	}
	for _, o := range ps {
		o(pdc)
	}
	if pdc.def == nil {
		pdc.err = errorx.New("Nil producer!")
	}
	return pdc
}

func (p *Producers) Send(c context.Context, message any) (err error) {

	if p.isGrayer.Judge(c) {
		return p.gray.Send(c, message)
	}
	if p.isShadow.Judge(c) {
		return p.shadow.Send(c, message)
	}

	return p.def.Send(c, message)
}

func (p *Producers) Error() (err error) {
	if p.def != nil && p.def.Error() != nil {
		return p.def.Error()
	}
	if p.shadow != nil && p.shadow.Error() != nil {
		return p.shadow.Error()
	}
	if p.gray != nil && p.gray.Error() != nil {
		return p.gray.Error()
	}
	return
}

func (p *Producers) Close() func() {
	return func() {
		if p.def != nil {
			p.def.Close()
		}
		if p.shadow != nil {
			p.shadow.Close()
		}
		if p.gray != nil {
			p.gray.Close()
		}
	}
}
