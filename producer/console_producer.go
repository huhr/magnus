package producer

import (
	"os"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/filter"
	"github.com/huhr/magnus/tools"
)

type ConsoleProducer struct {
	*BaseProducer
	reader *tools.UnitReader
}

func NewConsoleProducer(base *BaseProducer) (Producer, error) {
	return &ConsoleProducer{
		BaseProducer: base,
		reader:       tools.NewUnitReader(os.Stdin, base.config.Delimiter, base.config.BufSize),
	}, nil
}

func (cons *ConsoleProducer) Produce() {
	for cons.IsActive() {
		msg, err := cons.reader.ReadOne()
		if len(msg) > 0 {
			if !filter.Filter(msg, cons.config.Filters) {
				// ReadOne是阻塞的，停掉Produce时可能还会额外写入一个数据
				// console没有offset的概念，这条数据会被丢弃
				if cons.IsActive() {
					cons.pipe <- msg
				}
			}
		}
		if err != nil {
			log.Error(err.Error())
		}
	}
}
