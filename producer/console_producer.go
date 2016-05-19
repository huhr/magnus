package producer

import (
	"os"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/filter"
	"github.com/huhr/magnus/util"
)

// 输出到控制台
type ConsoleProducer struct {
	*BaseProducer
	reader *util.UnitReader
}

func NewConsoleProducer(base *BaseProducer) Producer {
	return &ConsoleProducer{
		BaseProducer: base,
		reader: util.NewUnitReader(os.Stdin, base.cfg.Delimiter, base.cfg.BufSize),
	}
}

func (cons *ConsoleProducer) Produce() {
	for cons.IsActive() {
		msg, err := cons.reader.ReadOne()
		if len(msg) > 0 {
			if !filter.Filter(msg, cons.cfg.Filters) {
				cons.pipe <- msg
			}
		}
		if err != nil {
			log.Error(err.Error())
		}
	}
}
