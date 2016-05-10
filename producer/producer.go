package producer

import (
	"github.com/huhr/magnus/config"
)

// Producer接口，负责生产数据
type Producer interface{
	// 生产一条数据
	Produce(chan []byte)
}

// 根据配置内容创建producer
func NewProducer(config.ProducerConfig) Producer {
	return nil
}

type BaseProducer struct {
}

