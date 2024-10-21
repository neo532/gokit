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
	name  string
	conf  *sarama.Config
	addrs []string
	group string
	err   error

	goCount          int
	bootstrapContext context.Context

	topics  []string
	handler *groupHandler

	consumer sarama.ConsumerGroup
}

func NewGroup(name string, addrs []string, group string, opts ...Option) (csm *ConsumerGroup) {

	// init parameter
	csm = &ConsumerGroup{
		name:    name,
		conf:    sarama.NewConfig(),
		addrs:   addrs,
		group:   group,
		goCount: runtime.NumCPU() / 2,
		handler: &groupHandler{
			name:     name,
			slowTime: 3 * time.Second,
			logger:   logger.NewDefaultILogger(),
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
	return csm.name
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
						queue.KeyName, csm.name,
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
	env        string
	name       string
	autoCommit bool
	handler    func(ctx context.Context, message []byte) (err error)
	slowTime   time.Duration
	logger     logger.ILogger
	msg        []byte
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *groupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *groupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {

	c := session.Context()
	var message []byte

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

	for msg := range claim.Messages() {

		message = msg.Value

		begin := time.Now()
		err := h.handler(c, message)
		cost := time.Since(begin)
		// mark ok
		session.MarkMessage(msg, "")
		// biz error
		if err != nil {
			h.logger.Error(c, "Handler has error!",
				queue.KeyName, h.name,
				queue.KeyErr, err,
				queue.KeyMessage, string(message),
			)
			continue
		}

		if !h.autoCommit {
			session.Commit()
		}

		// slow
		if cost > h.slowTime {
			h.logger.Warn(c, "slowlog",
				queue.KeyName, h.name,
				"slowTime", h.slowTime,
				"cost", cost,
				queue.KeyMessage, string(message),
			)
			continue
		}

		// if h.env == middleware.EnvProd && utf8.RuneCount(msg.Value) > log.MaxMsgLength {
		// 	msg.Value = []byte(string([]rune(string(msg.Value))[:log.MaxMsgLength]) + "...")
		// }
		h.logger.Info(c, string(message),
			queue.KeyName, h.name,
			queue.KeyPartition, msg.Partition,
			queue.KeyOffset, msg.Offset,
			"cost", cost,
		)
	}
	return nil
}
