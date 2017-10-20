package consumer

// 输出到控制台
type ThriftServiceConsumer struct {
	BaseConsumer
}

func NewThriftServiceConsumer(base BaseConsumer) (Consumer, error) {
	return &ThriftServiceConsumer{base}, nil
}

func (cons *ThriftServiceConsumer) Consume(msg []byte) bool {
	return true
}
