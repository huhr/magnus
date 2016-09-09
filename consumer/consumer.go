package consumer

import (
	"errors"

	"github.com/huhr/magnus/tools"
)

type Consumer interface {
	// 消费一条数据
	Consume([]byte) bool
	ShutDown()
}

type BaseConsumer struct {
	config tools.ConsumerConfig
	pipe   chan []byte
}

func (b BaseConsumer) ShutDown() {
	return
}

// 创建Consumer
func NewConsumer(config tools.ConsumerConfig, pipe chan []byte) (Consumer, error) {
	switch config.Consumer {
	case "console":
		return NewConsoleConsumer(BaseConsumer{config, pipe})
	case "archive":
		return NewArchiveConsumer(BaseConsumer{config, pipe})
	case "app":
		return NewAppConsumer(BaseConsumer{config, pipe})
	}
	return nil, errors.New("Illagle Consumer Type " + config.Consumer)
}
