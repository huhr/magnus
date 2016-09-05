// 对于数据流我们可能会有哪些应用场景？
// 磁盘文件 =》 MQ
// MQ =》磁盘、下游app
//
package stream

import (
	"errors"
	"sync"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/consumer"
	"github.com/huhr/magnus/producer"
	"github.com/huhr/magnus/tools"
)

// 多个Consumers获取流中的数据的方式
const (
	// 轮询
	ROUNDROBIN = iota + 1
	// 带权重随机
	WEIGHTEDRANDOM
)

// stream在业务上对应一条独立的数据流，每条数据流可以有多个
// producers和consumers，stream支持对数据按照一定策略进行分
// 发。
// 在多producers、consumers场景下，可能有不同的并发模式支持
// 详细设计下再写。
type Stream struct {
	// 管道用来缓存数据库
	Pipe      chan []byte
	Config    tools.StreamConfig
	consumers []consumer.Consumer
	producers []producer.Producer
}

func NewStream(config tools.StreamConfig) (*Stream, error) {
	stream := &Stream{
		Pipe:   make(chan []byte, config.CacheSize),
		Config: config,
	}
	if err := stream.initEnds(); err != nil {
		return nil, err
	}
	return stream, nil
}

// 创建stream对应的producers和consumers
func (s *Stream) initEnds() error {
	if len(s.Config.ProducerConfigs) == 0 || len(s.Config.ConsumerConfigs) == 0 {
		return errors.New("Producers or Consumers Config is nil")
	}
	// 初始化各个producers
	for _, config := range s.Config.ProducerConfigs {
		config.StreamName = s.Config.StreamName
		p, err := producer.NewProducer(config, s.Pipe)
		if err != nil {
			log.Error("%s Init Producer Error: %s", config.StreamName, err.Error())
			continue
		}
		s.producers = append(s.producers, p)
	}
	// 初始化各个consumer
	for _, config := range s.Config.ConsumerConfigs {
		config.StreamName = s.Config.StreamName
		c, err := consumer.NewConsumer(config, s.Pipe)
		if err != nil {
			log.Error("%s Init Consumer Error: %s", config.StreamName, err.Error())
			continue
		}
		s.consumers = append(s.consumers, c)
	}

	if len(s.producers) == 0 || len(s.consumers) == 0 {
		return errors.New("Deadends Stream")
	}
	log.Debug("Stream: %s, With: %d Producers, %d Consumers", s.Config.StreamName, len(s.producers), len(s.consumers))
	return nil
}

func (s *Stream) Run() {
	s.runSerial()
}

// 串行消费模式：保证数据顺序，先生产的producer先consume，数据逐一
// 进行consume，任意一个消费下游的阻塞都可能导致整个流的阻塞。可能
// 适用于处理有顺序需求的下游模块的场景。
// 这里Producer是并行的
func (s *Stream) runSerial() {
	var wg sync.WaitGroup
	// 启动各个生产协程
	for _, p := range s.producers {
		wg.Add(1)
		log.Debug("Producer %s Start", p.Name())
		go func(p producer.Producer) {
			defer func() {
				wg.Done()
				log.Debug("Producer %s Done", p.Name())
			}()
			p.Produce()
		}(p)
	}
	wg.Add(1)
	log.Debug("Consumers Start")
	go func() {
		defer func() {
			wg.Done()
			log.Debug("Consumers Done")
		}()
		s.Transit()
	}()

	wg.Wait()
	return
}

// 并行消费模式，多个consumer并行进行处理，不能保证数据的消费顺序，典型的
// 应用场景如下游为多个逻辑一直的回调发送模块
func (s *Stream) runConcurrent() {
	return
}

// 关闭程序，先关闭生产者，等数据消费完，记录offset
func (s *Stream) ShutDown() {
	for _, p := range s.producers {
		p.ShutDown()
	}
	close(s.Pipe)
}

// 根据不同的策略，将数据分发给不同的Consume
// 这里是穿行的，还没想好怎么调整为并发的
func (s *Stream) Transit() {
	var i int
	for msg := range s.Pipe {
		s.consumers[i].Consume(msg)
		i = (i + 1) % len(s.consumers)
	}
	// pipe已经关闭了，现在需要给所有的consumer发送一个EOF
	for _, _ = range s.consumers {
		continue
	}
}
