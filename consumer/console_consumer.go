package consumer

import (
	"fmt"

	"github.com/huhr/magnus/filter"
)

// 输出到控制台
type ConsoleConsumer struct {
	BaseConsumer
}

func NewConsoleConsumer(base BaseConsumer) (Consumer, error) {
	return &ConsoleConsumer{base}, nil
}

func (cons *ConsoleConsumer) Consume(msg []byte) bool {
	if !filter.Filter(msg, cons.config.Filters) {
		fmt.Printf("msg: %d, %s \n", len(msg), string(msg))
	}
	return true
}
