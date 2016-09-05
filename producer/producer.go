package producer

import (
	"errors"

	"github.com/huhr/magnus/tools"
)

// Producer接口，负责生产数据
type Producer interface {
	// 程序停机时，先停producer
	ShutDown()
	// 生产数据，produce函数是协程的执行体
	Produce()
	// 是否处于启动状态
	IsActive() bool
	// 获取当前producer的名称
	Name() string
}

// 根据配置内容创建producer
func NewProducer(config tools.ProducerConfig, pipe chan []byte) (Producer, error) {
	switch config.Producer {
	case "console":
		return NewConsoleProducer(NewBaseProducer(config, pipe))
	case "file":
		return NewFileProducer(NewBaseProducer(config, pipe))
	}
	return nil, errors.New("Illagle Producer Type " + config.Producer)
}

type BaseProducer struct {
	config tools.ProducerConfig
	pipe   chan []byte
	isOff  bool
}

func NewBaseProducer(config tools.ProducerConfig, pipe chan []byte) *BaseProducer {
	base := &BaseProducer{config: config, pipe: pipe}
	return base
}

func (base *BaseProducer) ShutDown() {
	base.isOff = true
}

func (base BaseProducer) IsActive() bool {
	return !base.isOff
}

func (base BaseProducer) Name() string {
	return base.config.WorkerName
}
