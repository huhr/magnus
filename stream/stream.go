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
	// 广播
	BROADCAST
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
	s.run1()
}

// 第一种模式，生产者并发产生数据，数据在由stream按照调度策略逐一发送给各个
// consumers，发送失败顺延为下一次发送
func (s *Stream) run1() {
	var wg sync.WaitGroup
	// 启动各个生产协程，这里是没有顺序关系的
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

func (s *Stream) runConcurrent() {
	return
}

// 关闭程序，先关闭生产者，等数据消费完，记录offset
func (s *Stream) ShutDown() {
	log.Debug("Stream %s Begin To ShutDown", s.Config.StreamName)
	for _, p := range s.producers {
		p.ShutDown()
	}
	close(s.Pipe)
}

// 根据不同的策略，将数据分发给不同的Consume
func (s *Stream) Transit() {
	switch s.Config.TransitType {
	case ROUNDROBIN:
		s.transitByROUNDROBIN()
	case BROADCAST:
		s.transitByBROADCAST()
	case WEIGHTEDRANDOM:
		s.transitByWEIGHTEDRANDOM()
	default:
		log.Error("Undefined TransitType %d", s.Config.TransitType)
	}
}

// 顺序逐个发送，发送失败后顺延发送给下一个consumer，默认一个数据三轮发送失败
// 后会被丢弃，连续三个数据三轮发送失败，Stream停止运行
func (s *Stream) transitByROUNDROBIN() {
	var i = -1                               // consumers的下标
	var errMsgs = 0                          // 连续消费出错的数据条数
	var roundRetryNum = 2 * len(s.consumers) // 当一条消息消费出错时，顺延消费两圈
	for msg := range s.Pipe {
		j := 0
		for j < roundRetryNum {
			j++
			i = (i + 1) % len(s.consumers)
			if s.consumers[i].Consume(msg) {
				break
			}
		}
		if j >= roundRetryNum {
			log.Error("Stream %s retry too much numbers, msg: %s", s.Config.StreamName, msg)
			errMsgs++
			if errMsgs > 2 {
				log.Error("Stream %s has occurred exception for 3 messages", s.Config.StreamName)
				s.ShutDown()
			}
		} else {
			errMsgs = 0
		}
	}
	for _, c := range s.consumers {
		c.ShutDown()
	}
}

// 广播模式，给每一个consumer发送数据，考虑阻塞的形式以及重试机制
func (s *Stream) transitByBROADCAST() {
	for msg := range s.Pipe {
		for _, c := range s.consumers {
			c.Consume(msg)
		}
	}
	for _, c := range s.consumers {
		c.ShutDown()
	}
}

// 带权重抽取
func (s *Stream) transitByWEIGHTEDRANDOM() {
	return
}
