package producer

import (
	"os"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/config"
	"github.com/huhr/magnus/filter"
)

// 输出到控制台
type ConsoleProducer struct {
	BaseProducer
	reader *UnitReader
}

func NewConsoleProducer(cfg config.ProducerConfig, pipe chan []byte) *ConsoleProducer {
	base := BaseProducer{cfg: cfg, pipe: pipe}
	return &ConsoleProducer{
		BaseProducer: base,
		reader: NewUnitReader(os.Stdin, cfg.Delimiter, cfg.BufSiZe),
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
