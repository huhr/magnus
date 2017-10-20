package consumer

import "fmt"

// 输出到http服务
type HttpServiceConsumer struct {
	BaseConsumer
}

func NewHttpServiceConsumer(base BaseConsumer) (Consumer, error) {
	return &HttpServiceConsumer{base}, nil
}

func (cons *HttpServiceConsumer) Consume(msg []byte) bool {
	fmt.Printf("msg: %d, %s \n", len(msg), string(msg))
	return true
}
