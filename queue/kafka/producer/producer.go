package producer

/*
 * @abstract producer
 * @mail neo532@126.com
 * @date 2024-10-20
 */

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
)

var (
	_ queue.Producer = (*Producer)(nil)

	instanceLock sync.Mutex
	producerMap  = make(map[string]*Producer, 2)
)

type EncodeMessageFunc func(message interface{}) (msg []byte, err error)

func JsonMessageEncoder(message interface{}) (msg []byte, err error) {
	return json.Marshal(message)
}

type Producer struct {
	Name    string         `json:"name"`
	Conf    *sarama.Config `json:"config"`
	Addrs   []string       `json:"name"`
	IsAsync bool           `json:"is_async"`
	Topic   string         `json:"topic"`

	syncProducer  sarama.SyncProducer  `json:"-"`
	asyncProducer sarama.AsyncProducer `json:"-"`

	key              string                     `json:"-"`
	logger           logger.ILogger             `json:"-"`
	encoder          EncodeMessageFunc          `json:"-"`
	close            func()                     `json:"-"`
	err              error                      `json:"-"`
	bootstrapContext context.Context            `json:"-"`
	middleware       []queue.ProducerMiddleware `json:"-"`
}

func New(name string, addrs []string, opts ...Option) (pdc *Producer) {
	// one instance
	instanceLock.Lock()
	defer instanceLock.Unlock()

	// init
	pdc = &Producer{
		Name:             name,
		Conf:             sarama.NewConfig(),
		Addrs:            addrs,
		logger:           logger.NewDefaultILogger(),
		bootstrapContext: context.Background(),
		encoder:          JsonMessageEncoder,
		middleware:       make([]queue.ProducerMiddleware, 0, 1),
	}
	pdc.Conf.Version = sarama.V0_11_0_2
	pdc.Conf.Producer.Return.Successes = true
	for _, o := range opts {
		o(pdc)
	}

	if b, e := json.Marshal(pdc); e == nil {
		pdc.key = fmt.Sprintf("%x", md5.Sum(b))
	}

	if p, ok := producerMap[pdc.key]; ok {
		pdc = p
		return
	}

	ps := []interface{}{
		queue.KeyName, pdc.Name,
		queue.KeyIsAsync, pdc.IsAsync,
	}

	// validate
	if pdc.err = pdc.Conf.Validate(); pdc.err != nil {
		ps = append(ps,
			queue.KeyConfig, pdc.Conf,
			queue.KeyErr, pdc.err,
		)
		pdc.logger.Error(pdc.bootstrapContext, "conf.Validate Has error.", ps...)
		return
	}

	switch pdc.IsAsync {
	case false:
		pdc.syncProducer, pdc.err = sarama.NewSyncProducer(pdc.Addrs, pdc.Conf)
		pdc.close = func() {
			if pdc.syncProducer != nil {
				pdc.err = pdc.syncProducer.Close()
			}
		}
	case true:
		pdc.asyncProducer, pdc.err = sarama.NewAsyncProducer(pdc.Addrs, pdc.Conf)
		pdc.close = func() {
			if pdc.asyncProducer != nil {
				pdc.err = pdc.asyncProducer.Close()
			}
		}
		go func() {
			for {
				select {
				case e := <-pdc.asyncProducer.Errors():
					if e != nil {
						b, _ := e.Msg.Value.Encode()
						ps = append(ps,
							queue.KeyErr, e.Error(),
							queue.KeyTopic, e.Msg.Topic,
							queue.KeyOffset, e.Msg.Offset,
							queue.KeyPartition, e.Msg.Partition,
							queue.KeyKey, e.Msg.Key,
							queue.KeyValue, string(b),
						)
						pdc.logger.Error(pdc.bootstrapContext, "Async producer has error!", ps...)
					}
				case <-pdc.asyncProducer.Successes():
				case <-pdc.bootstrapContext.Done():
					return
				}
			}
		}()
	}
	if pdc.err != nil {
		ps = append(ps, queue.KeyErr, pdc.err)
		pdc.logger.Error(pdc.bootstrapContext, "sarama.NewProducer has error!", ps...)
		return
	}
	pdc.logger.Info(pdc.bootstrapContext, "Producer Running!", ps...)
	return
}

func (pdc *Producer) Err() error {
	return pdc.err
}

func (pdc *Producer) Send(c context.Context, message interface{}) (err error) {

	c = queue.InitHeaderToContext(c)

	ps := []interface{}{
		queue.KeyName, pdc.Name,
		queue.KeyIsAsync, pdc.IsAsync,
	}

	var msg []byte
	if msg, err = pdc.encoder(message); err != nil {
		ps = append(ps, queue.KeyErr, err, queue.KeyMessage, message)
		pdc.logger.Error(c, "Producer's encoder Has err!", ps...)
		return
	}
	ps = append(ps, queue.KeyMessage, string(msg))

	h := func(c context.Context, message interface{}) (err error) {

		pm := &sarama.ProducerMessage{
			Topic:     pdc.Topic,
			Timestamp: time.Now(),
			Value:     sarama.ByteEncoder(msg),
			Headers:   []sarama.RecordHeader{}, // at leaset kafka v0.11+
		}
		if h, ok := queue.GetHeaderFromContext(c); ok {

			if hk := h.Value(queue.KeyHashKey); hk != "" {
				pm.Key = sarama.StringEncoder(hk)
				ps = append(ps, queue.KeyHashKey, hk)
			}

			h.Range(func(k string, v string) bool {
				pm.Headers = append(pm.Headers, sarama.RecordHeader{
					[]byte(k), []byte(v),
				})
				return true
			})
		}

		switch pdc.IsAsync {
		case false:
			var p int32
			var o int64
			p, o, err = pdc.syncProducer.SendMessage(pm)
			ps = append(ps,
				queue.KeyPartition, p,
				queue.KeyOffset, o,
			)
		case true:
			pdc.asyncProducer.Input() <- pm
		}
		return
	}

	if len(pdc.middleware) > 0 {
		h = queue.ChainProducer(pdc.middleware...)(h)
	}

	if err = h(c, message); err != nil {
		ps = append(ps, queue.KeyErr, err)
		pdc.logger.Error(c, "Producer's sending Has err!", ps...)
		return
	}

	pdc.logger.Info(c, "Producer have been delivered!", ps...)
	return
}

func (pdc *Producer) Close() func() {
	return pdc.close
}
