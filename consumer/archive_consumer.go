package consumer

import (
	"fmt"
	"io"
	"os"

	log "github.com/huhr/simplelog"
)

// 写到文件中
type ArchiveConsumer struct {
	BaseConsumer
	writer io.Writer
}

func NewArchiveConsumer(base BaseConsumer) (Consumer, error) {
	var w io.Writer
	if _, err := os.Stat(base.config.FilePath); os.IsNotExist(err) {
		w, err = os.Create(base.config.FilePath)
		if err != nil {
			return nil, err
		}
	} else {
		w, err = os.OpenFile(base.config.FilePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
	}
	archive := &ArchiveConsumer{base, w}
	return archive, nil
}

func (f *ArchiveConsumer) Consume(msg []byte) bool {
	outPut := fmt.Sprintf("%s\n", msg)
	_, err := f.writer.Write([]byte(outPut))
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}
