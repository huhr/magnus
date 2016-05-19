package producer

import (
	"fmt"
	"syscall"
	"os"
	"time"
	//"path/filepath"
	"io"
	"io/ioutil"
	"errors"
	"encoding/json"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/filter"
	"github.com/huhr/magnus/util"
)


// 输出到控制台
type FileProducer struct {
	*BaseProducer
	offset          int64
	// 当前正在处理的文件
	currentFile     *os.File
	reader			*util.UnitReader
}

// 创建FileProducer实例，尝试去加载runtime文件，找到需要
// 打开的文件，seek到上次读的位置，继续读
func NewFileProducer(base *BaseProducer) Producer {
	producer := &FileProducer{
		BaseProducer: base,
	}
	// 加载运行时文件出错
	if producer.loadRuntime() != nil {
		return nil
	}
	return producer
}

func (f *FileProducer) Produce() {
	for f.IsActive() {
		// read一个数据单元
		msg, err := f.reader.ReadOne()
		if len(msg) > 0 {
			// read到数据了，记录offset
			f.offset += int64(len(msg)) + int64(len(f.cfg.Delimiter))
			// 过滤数据，不符合规则的丢弃
			if !filter.Filter(msg, f.cfg.Filters) {
				f.pipe <- msg
			}
			continue
		}
		// 没有读到数据，读到文件尾，判断是不是需要进行文件切换
		if err == io.EOF {
			// 检查当前的文件和配置的路径文件是不是同一个文件
			f1, _ := os.Stat(f.cfg.FilePath)
			f2, _ := f.currentFile.Stat()
			// 没有发生切割，继续读数据
			if os.SameFile(f1, f2) {
				continue
			} else {
				// 已经发生了切割，open新的文件直到成功
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
	f.dumpRuntime()
	close(f.pipe)
}

type runtimeData struct {
	Offset      int64
	Dev			uint64
	Ino         uint64
}

// 程序启动时加载运行时文件，找到上次读的文件
func (f *FileProducer) loadRuntime() error {
	// 找到了对应文件
	if runtimeFile, err := os.Open(fmt.Sprintf("runtime/stream_%s_%s.data",
		f.cfg.StreamName, f.cfg.WorkerName)); err == nil {
		var data runtimeData
		buf, err := ioutil.ReadAll(runtimeFile)
		if err != nil {
			return f.seekFile(nil)
		}
		if err = json.Unmarshal(buf, &data); err != nil {
			log.Error("json unmarshal runtime file error in %s-%s, err: %s",
				f.cfg.StreamName, f.cfg.WorkerName, err.Error())
			return f.seekFile(nil)
		}
		return f.seekFile(&data)
	}

	return f.seekFile(nil)
}

// 找到需要打开的文件，并seek到指定的offset
func (f *FileProducer) seekFile(data *runtimeData) error {
	var file *os.File
	var err error
	// 没有runtime文件，直接打开filePath对应的文件
	if data == nil {
		if file, err = os.Open(f.cfg.FilePath); err != nil {
			log.Error("open file error: %s", err.Error())
			return err
		}
		f.reader = util.NewUnitReader(file, f.cfg.Delimiter, f.cfg.BufSize)
		f.currentFile = file
		return nil
	}
	// 先打开，再获取stat，这个操作反过来由于不是原子操作，可能会导致异常
	if file, err = os.Open(f.cfg.FilePath); err == nil {
		fileInfo, _ := file.Stat()
		// 当前文件正是要继续读的文件
		if fileInfo.Sys().(*syscall.Stat_t).Dev == data.Dev && fileInfo.Sys().(*syscall.Stat_t).Ino == data.Ino {
			if _, err := file.Seek(data.Offset, 0); err != nil {
				log.Error("seek error: %s", err.Error())
				return err
			}
			f.offset = data.Offset
			f.reader = util.NewUnitReader(file, f.cfg.Delimiter, f.cfg.BufSize)
			f.currentFile = file
			return nil
		}
		if file, err = os.Open(f.cfg.FilePath); err != nil {
			log.Error("open file error: %s", err.Error())
			return err
		}
		f.reader = util.NewUnitReader(file, f.cfg.Delimiter, f.cfg.BufSize)
		f.currentFile = file
		return nil
	}
	return errors.New("还没写好")
}

// 停止程序时，需要dump运行时状态到文件中
func (f FileProducer) dumpRuntime() {
	runtimeFile, err := os.Create(fmt.Sprintf("runtime/stream_%s_%s.data", f.cfg.StreamName, f.cfg.WorkerName))
	if err != nil && !os.IsExist(err) {
		panic("create data file error:" + err.Error())
	}
	data := runtimeData{}
	stat, _ := f.currentFile.Stat()
	data.Offset = f.offset
	data.Dev = stat.Sys().(*syscall.Stat_t).Dev
	data.Ino = stat.Sys().(*syscall.Stat_t).Ino

	dataJson, _ := json.Marshal(data)
	runtimeFile.Write(dataJson)
}

