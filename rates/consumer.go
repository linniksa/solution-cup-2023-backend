package rates

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"solution/config"
	"solution/ent/schema"
	"strconv"
	"sync"
	"time"
)

type Manager struct {
	reader *kafka.Reader
	logger *zap.Logger
	m      sync.RWMutex
	data   *ExchangeInfo
}

func New(
	cfg *config.Config,
	logger *zap.Logger,
) (*Manager, error) {
	host := cfg.KafkaHost + ":" + strconv.Itoa(cfg.KafkaPort)

	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{host},
		Topic:   cfg.KafkaTopicName,
		GroupID: "my-group",
	})

	s := &Manager{
		reader: r,
		logger: logger,
	}

	ch := make(chan struct{}, 1)
	go s.worker(ch)
	go func() {
		time.Sleep(1 * time.Second)
		ch <- struct{}{}
	}()

	<-ch

	return s, nil
}

func (s *Manager) worker(ch chan struct{}) {
	ctx := context.Background()

	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			s.logger.Error("error reading message", zap.Error(err))
			continue
		}
		s.convertRateByte(msg.Value)
		select {
		case ch <- struct{}{}:
			// noop
		default:
		}
	}
}

func (s *Manager) convertRateByte(value []byte) {
	s.m.Lock()
	defer s.m.Unlock()

	var rate map[string]float64
	if err := json.Unmarshal(value, &rate); err != nil {
		s.logger.Error("Wrong format?", zap.Error(err))
		return
	}

	s.convertRate(rate)
}

func (s *Manager) convertRate(rate map[string]float64) {
	info := &ExchangeInfo{
		data: make([]exchangeInfo, 0, 10),
	}

	for c, v := range rate {
		i, ok := toInfo(s.logger, c, v)
		if !ok {
			continue
		}

		info.data = append(info.data, i)
	}

	s.data = info
}

func (s *Manager) Rate() *ExchangeInfo {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.data
}

func toInfo(logger *zap.Logger, c string, p float64) (exchangeInfo, bool) {
	if len(c) != 6 {
		logger.Error("wrong currency?", zap.String("currency", c))
		return exchangeInfo{}, false
	}

	from := schema.Currency(c[:3])
	to := schema.Currency(c[3:6])

	if !from.Valid() || !to.Valid() {
		logger.Error("wrong currencies?", zap.String("currency", c))
		return exchangeInfo{}, false
	}

	return exchangeInfo{
		From: from,
		To:   to,
		Rate: p,
	}, true

}
