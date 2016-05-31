package stream

import (
	"io"
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
	ROUNDROBIN = iota + 1
	// 广播
	BROADCAST
	// 带权重随机
	WEIGHTEDRANDOM
)

// 负责缓存中转数据，每个stream是一个独立的数据流，
// 每条数据流可以对应多个producers和consumers
type Stream struct{
	Name    string
	TransitType    int
	Pipe	chan []byte
	Cfg     config.StreamConfig
	consumers []consumer.Consumer
	producers []producer.Producer
}

func NewStream(cfg config.StreamConfig) *Stream {
	return &Stream{
		Name: cfg.Name,
		TransitType: cfg.TransitType,
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
		cfg.StreamName = s.Name
		p := producer.NewProducer(cfg, s.Pipe)
		if p == nil {
			continue
		}
		s.producers = append(s.producers, p)
	}
	for _, cfg := range s.Cfg.Ccfgs {
		cfg.StreamName = s.Name
		c := consumer.NewConsumer(cfg, s.Pipe)
		if c == nil {
			continue
		}
		s.consumers = append(s.consumers, consumer.NewConsumer(cfg, s.Pipe))
	}
	return nil
}

func (s *Stream) Run() {
	var wg sync.WaitGroup
	if err := s.initEnds(); err != nil {
		log.Error("Init stream %s fail: %s", s.Name, err.Error())
		return
	}
	log.Debug("stream: %s, total: %d producers, %d consumers", s.Name, len(s.producers), len(s.consumers))
	// 启动各个生产协程
	for _, p := range s.producers {
		wg.Add(1)
		log.Debug("producer %s start", p.Name())
		go func(p producer.Producer) {
			defer func() {
				wg.Done()
				log.Debug("producer %s done", p.Name())
			}()
			p.Produce()
		}(p)
	}
	// 启动消费协程
	wg.Add(1)
	log.Debug("consumer start")
	go func() {
		defer func() {
			wg.Done()
			log.Debug("consumer start")
		}()
		s.Transit()
	}()

	wg.Wait()
	return
}

func (s *Stream) ShutDown() {
	for _, p := range s.producers {
		p.ShutDown()
	}
	close(s.Pipe)
}


// 根据不同的策略，将数据分发给不同的Consume
// 这里需要是并发的，但是要限制一定的并发度
func (s *Stream) Transit() {
	var i int
	for msg := range s.Pipe {
		s.consumers[i].Consume(msg)
		i = (i + 1) % len(s.consumers)
	}
	// pipe已经关闭了，现在需要给所有的consumer发送一个EOF
	for _, c := range s.consumers {
		c.Consume()
	}
}

