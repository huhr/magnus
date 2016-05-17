package producer

import (
	"fmt"
	"os"
	"time"
	"io"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/config"
	"github.com/huhr/magnus/filter"
	"github.com/huhr/magnus/util"
)

// 输出到控制台
type FileProducer struct {
	BaseProducer
	currentFile		*os.File
	reader			*util.UnitReader
}

func NewFileProducer(cfg config.ProducerConfig, pipe chan []byte) *FileProducer {
	var file *os.File
	var err error

	if file, err = os.Open(cfg.FilePath); err != nil && os.IsNotExist(err) {
		panic(fmt.Sprintf("%s is not exist", cfg.FilePath))
	}
	return &FileProducer{
		BaseProducer: BaseProducer{cfg: cfg, pipe: pipe},
		currentFile: file,
		reader: util.NewUnitReader(file, cfg.Delimiter, cfg.BufSize),
	}
}

func (f *FileProducer) Produce() {
	for f.IsActive() {
		msg, err := f.reader.ReadOne()
		if len(msg) > 0 {
			if !filter.Filter(msg, f.cfg.Filters) {
				f.pipe <- msg
			}
		}
		// 读到文件尾，判断是不是需要进行文件切换
		if len(msg) == 0 && err == io.EOF {
			f1, _ := os.Stat(f.cfg.FilePath)
			f2, _ := f.currentFile.Stat()
			if os.SameFile(f1, f2) {
				continue
			} else {
				for true {
					if f.currentFile, err = os.Open(f.cfg.FilePath); err != nil {
						log.Error("File Roll: open file %s error: %s", f.cfg.FilePath, err.Error())
						time.Sleep(3 * time.Second)
						continue
					}
					f.reader.ResetReader(f.currentFile)
					log.Debug("Roll File Success")
					break
				}
				continue
			}
		}
	}
}
