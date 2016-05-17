package stream

import (
	"errors"
	"sync"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/config"
	"github.com/huhr/magnus/consumer"
	"github.com/huhr/magnus/producer"
)

// 流转方式
const (
	// 轮询
	ROUNDROBIN = iota
	// 广播
	BROADCAST
	// 带权重随机
	WEIGHTEDRANDOM
)

// 负责缓存中转数据，每个stream是一个独立的数据流，
// 每条数据流可以对应多个producers和consumers
type Stream struct{
	Name    string
	RollType    int
	Pipe	chan []byte
	Cfg     config.StreamConfig
	consumers []consumer.Consumer
	producers []producer.Producer
}

func NewStream(cfg config.StreamConfig) *Stream {
	return &Stream{
		Name: cfg.Name,
		RollType: cfg.RollType,
		Pipe: make(chan []byte, cfg.CacheSize),
		Cfg: cfg,
	}
}

// 创建stream两端的生产消费对象
func (s *Stream) initEnds() error {
	if len(s.Cfg.Pcfgs) == 0 || len(s.Cfg.Ccfgs) == 0 {
		log.Error("producer or consumer is missing")
		return errors.New("producer or consumer is missing")
	}
	for _, cfg := range s.Cfg.Pcfgs {
		s.producers = append(s.producers, producer.NewProducer(cfg, s.Pipe))
	}
	for _, cfg := range s.Cfg.Ccfgs {
		s.consumers = append(s.consumers, consumer.NewConsumer(cfg, s.Pipe))
	}
	return nil
}

func (s *Stream) Run() {
	var wg sync.WaitGroup
	log.Debug("stream init")
	if s.initEnds() != nil {
		return
	}
	log.Debug("stream: %s, total: %d producers", s.Name, len(s.producers))
	for _, p := range s.producers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Produce()
		}()
	}
	go func() {
		s.Roll()
	}()

	wg.Wait()
	return
}

func (s *Stream) Roll() {
	for true {
		for _, c := range s.consumers {
			c.Consume()
		}
	}
}

