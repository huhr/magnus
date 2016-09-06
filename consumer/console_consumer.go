package consumer

import "fmt"

// 输出到控制台
type ConsoleConsumer struct {
	BaseConsumer
}

func NewConsoleConsumer(base BaseConsumer) (Consumer, error) {
	return &ConsoleConsumer{base}, nil
}

func (cons *ConsoleConsumer) Consume(msg []byte) bool {
	fmt.Printf("msg: %d, %s \n", len(msg), string(msg))
	return true
}
