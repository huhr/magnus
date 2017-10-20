// kafka producer使用github.com/Shopify/sarama作为客户端来链接Kafka
//
//
//
package producer

import (
	"fmt"
	"os"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/filter"
	"github.com/huhr/magnus/tools"
)

type KafkaProducer struct {
	*BaseProducer
	reader *tools.UnitReader
}

func NewKafkaProducer(base *BaseProducer) (Producer, error) {
	kumer, err := NewConsumer([]string{"10.0.0.206:10093"}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return &KafkaProducer{
		BaseProducer: base,
		reader:       tools.NewUnitReader(os.Stdin, base.config.Delimiter, base.config.BufSize),
	}, nil
}

func (k *KafkaProducer) Produce() {
	for k.IsActive() {
		msg, err := k.reader.ReadOne()
		if len(msg) > 0 {
			if !filter.Filter(msg, k.config.Filters) {
				k.pipe <- msg
			}
		}
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (k *KafkaProducer) ShutDown() {
	k.Consumer.Close()
}
