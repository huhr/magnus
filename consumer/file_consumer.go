package consumer

import (
	"fmt"

	"github.com/huhr/magnus/filter"
)

// 写到文件中
type FileConsumer struct {
	BaseConsumer
}

func NewFileConsumer(base BaseConsumer) (Consumer, error) {
	return &FileConsumer{base}, nil
}

func (f *FileConsumer) Consume(msg []byte) bool {
	if !filter.Filter(msg, f.config.Filters) {
		fmt.Printf("msg: %d, %s \n", len(msg), string(msg))
	}
	return true
}
