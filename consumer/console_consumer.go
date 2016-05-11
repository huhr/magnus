package consumer

import (
	"fmt"

	"github.com/huhr/magnus/filter"
)

// 输出到控制台
type ConsoleConsumer struct {
	BaseConsumer
}

func (cons *ConsoleConsumer) Consume() {
	msg := <-cons.pipe
	if !filter.Filter(msg, cons.cfg.Filters) {
		fmt.Printf("%s\n", string(msg))
		return
	}
	// golang没有对尾递归进行优化
	cons.Consume()
}
