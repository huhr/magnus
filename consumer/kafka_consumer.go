// 先尝试下AsyncProducer
package consumer

import "fmt"

// 输出到控制台
type KafkaConsumer struct {
	BaseConsumer
}

func NewKafkaConsumer(base BaseConsumer) (Consumer, error) {
	return &KafkaConsumer{base}, nil
}

func (cons *KafkaConsumer) Consume(msg []byte) bool {
	fmt.Printf("msg: %d, %s \n", len(msg), string(msg))
	return true
}
