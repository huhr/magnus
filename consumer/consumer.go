package consumer

import (
	"github.com/huhr/magnus/config"
)

type Consumer interface {
	// 消费一条数据
	Consume([]byte) bool
}

type BaseConsumer struct {
	cfg		config.ConsumerConfig
	pipe	chan []byte
}

// 创建Consumer
func NewConsumer(cfg config.ConsumerConfig, pipe chan []byte) Consumer {
	switch cfg.Consumer {
	case "console":
		return &ConsoleConsumer{BaseConsumer{cfg, pipe}}
	case "app":
		return NewAppConsumer(BaseConsumer{cfg, pipe})
	}
	return nil
}
