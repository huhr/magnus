package consumer

import (
	"github.com/huhr/magnus/config"
)

//
type Consumer interface{
	// 消费一条数据
	Consume(chan []byte)
}

// 根据配置内容创建consumer
func NewConsumer(cfg config.ConsumerConfig) Consumer {
	return nil
}
