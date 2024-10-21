package consumergroup

/*
 * @abstract consumer
 * @mail neo532@126.com
 * @date 2024-10-21
 */

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/IBM/sarama"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
)

var _ queue.Consumer = (*ConsumerGroup)(nil)

// ========== ConsumerGroup ==========
type ConsumerGroup struct {
	conf    *sarama.Config
	addrs   []string
	group   string
	topics  []string
	handler *groupHandler

	err              error
	goCount          int
	bootstrapContext context.Context

	consumer sarama.ConsumerGroup
}

func NewGroup(name string, addrs []string, group string, opts ...Option) (csm *ConsumerGroup) {

	// init parameter
	csm = &ConsumerGroup{
		conf:    sarama.NewConfig(),
		addrs:   addrs,
		group:   group,
		goCount: runtime.NumCPU() / 2,
		handler: &groupHandler{
			name:       name,
			slowTime:   3 * time.Second,
			logger:     logger.NewDefaultILogger(),
			middleware: make([]queue.ConsumerMiddleware, 0, 1),
		},
		bootstrapContext: context.Background(),
	}
	if csm.goCount < 3 {
		csm.goCount = 3
	}
	csm.conf.Version = sarama.V0_11_0_2
	csm.handler.autoCommit = csm.conf.Consumer.Offsets.AutoCommit.Enable
	csm.conf.Consumer.MaxWaitTime = time.Second
	for _, o := range opts {
		o(csm)
	}

	// check
	if csm.err = csm.conf.Validate(); csm.err != nil {
		csm.handler.logger.Error(csm.bootstrapContext, "Validate has error",
			queue.KeyConfig, csm.conf,
			queue.KeyErr, csm.err,
		)
		return
	}

	// initilize
	if csm.consumer, csm.err = sarama.NewConsumerGroup(
		csm.addrs,
		csm.group,
		csm.conf); csm.err != nil {
		csm.handler.logger.Error(csm.bootstrapContext, "NewGroup has error!",
			queue.KeyErr, csm.err,
		)
		return
	}
	return
}

func (csm *ConsumerGroup) Name() (name string) {
	return csm.handler.name
}

func (csm *ConsumerGroup) Stop(c context.Context) (err error) {
	if csm.consumer != nil {
		err = csm.consumer.Close()
	}
	return
}

func (csm *ConsumerGroup) Start(c context.Context) (err error) {
	for i := 0; i < csm.goCount; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(runtime.Caller(0))
					csm.handler.logger.Error(c, "Start has panic!",
						queue.KeyErr, err,
						"track", string(debug.Stack()),
					)
				}
			}()

			for {
				select {
				case <-c.Done():
					csm.handler.logger.Info(c, "topic consumer have canceled!",
						queue.KeyTopic, csm.topics,
					)
					return
				default:
					csm.handler.logger.Info(c, "Consumer is starting!",
						queue.KeyName, csm.handler.name,
						queue.KeyTopic, strings.Join(csm.topics, ","),
						queue.KeyAddr, strings.Join(csm.addrs, ","),
						queue.KeyGroup, csm.group,
					)

					// This method blocks until the consumer service is stopped
					if err := csm.consumer.Consume(c, csm.topics, csm.handler); err != nil {
						csm.handler.logger.Error(c, "Consume has error",
							queue.KeyErr, err,
						)
						return
					}
				}
			}
		}()
	}
	return
}

// Consumer represents a Sarama consumer group consumer
type groupHandler struct {
	name       string
	autoCommit bool
	handler    func(ctx context.Context, message []byte) (err error)
	slowTime   time.Duration
	logger     logger.ILogger
	middleware []queue.ConsumerMiddleware
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *groupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *groupHandler) Cleanup(session sarama.ConsumerGroupSession) (err error) {
	return
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {

	c := session.Context()
	var message []byte
	var ps []interface{}

	defer func() {
		if err := recover(); err != nil {
			h.logger.Error(c, "Handler has panic!",
				queue.KeyErr, err,
				queue.KeyMessage, string(message),
				queue.KeyStack, string(debug.Stack()),
			)
		}
		time.Sleep(time.Second)
	}()

	hdl := func(c context.Context, message []byte) (err error) {
		err = h.handler(c, message)
		return
	}

	for m := range claim.Messages() {

		message = m.Value
		ps = []interface{}{
			queue.KeyName, h.name,
			queue.KeyPartition, m.Partition,
			queue.KeyOffset, m.Offset,
			queue.KeyMessage, string(message),
		}

		c = queue.InitHeaderToContext(c)
		if header, ok := queue.GetHeaderFromContext(c); ok {
			for _, h := range m.Headers {
				header.Set(string(h.Key), string(h.Value))
			}
		}

		if len(h.middleware) > 0 {
			hdl = queue.ChainConsumer(h.middleware...)(hdl)
		}

		begin := time.Now()
		if err = hdl(c, message); err != nil {
			ps = append(ps, queue.KeyErr, err)
			h.logger.Error(c, "Consumer's Has err!", ps...)
			return
		}
		cost := time.Since(begin)
		ps = append(ps, "cost", cost)

		// mark ok
		session.MarkMessage(m, "")

		// biz error
		if err != nil {
			ps = append(ps, queue.KeyErr, err)
			h.logger.Error(c, "Handler has error!", ps...)
			continue
		}

		if !h.autoCommit {
			session.Commit()
		}

		// slow
		if cost > h.slowTime {
			ps = append(ps,
				"slowTime", h.slowTime,
			)
			h.logger.Warn(c, "slowlog", ps...)
			continue
		}

		h.logger.Info(c, "", ps...)
		return
	}

	return
}
