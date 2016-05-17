package producer

import (
	"github.com/huhr/magnus/config"
)

// Producer接口，负责生产数据
type Producer interface{
	// 程序停机时，先停producer
	ShutDown()
	// 生产数据，produce函数是协程的执行体
	Produce()
	// 是否处于启动状态
	IsActive() bool
}

// 根据配置内容创建producer
func NewProducer(cfg config.ProducerConfig, pipe chan []byte) Producer {
	switch cfg.Producer {
	case "console":
		return NewConsoleProducer(cfg, pipe)
	case "file":
		return NewFileProducer(cfg, pipe)
	}
	return nil
}

type BaseProducer struct {
	cfg config.ProducerConfig
	pipe chan []byte
	isOff bool
}

func (base *BaseProducer) ShutDown() {
	base.isOff = true
}

func (base *BaseProducer) IsActive() bool {
	return !base.isOff
}
