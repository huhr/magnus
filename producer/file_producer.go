package producer

import (
	"fmt"
	"syscall"
	"strconv"
	"os"
	"time"
	"path/filepath"
	"regexp"
	"io"
	"io/ioutil"
	"errors"
	"encoding/json"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/filter"
	"github.com/huhr/magnus/util"
)

const (
	PRESENT = iota + 1
	ORIGINAL
)

// 这里实现了日志轮转切割的功能，通过固定规则的文件名和inode
// 信息来校验文件，处理的文件需要保证不能使用vi进行copy & mv
// 类型的编辑保存操作
type FileProducer struct {
	*BaseProducer
	offset          int64
	// 当前正在处理的文件
	runtime         *runtimeData
	reader			*util.UnitReader
}

// 创建FileProducer实例时，尝试去加载runtime文件，找到需要
// 打开的文件，seek到上次读的位置，继续读
func NewFileProducer(base *BaseProducer) Producer {
	producer := &FileProducer{
		BaseProducer: base,
	}
	if producer.loadRuntime() != nil {
		return nil
	}
	if producer.seekFile() != nil {
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
		// 读到EOF:
		//		1、当前文件在backup文件夹：seek下一个文件
		//      2、当前文件还在working目录：sleep
		if err == io.EOF {
			// 可能是已经mv了旧文件，新的文件还没有生成，打错误日志
			same, err := f.match(f.cfg.FilePath)
			// 路径出错，working目录的文件还没有生成，当前读的文件一定不再working目录了
			// 可以尝试去寻找下一个要读的文件继续读，working目录没有文件
			if same == false || err != nil {
				if err != nil && !os.IsNotExist(err) {
					log.Error("%s-%s file roll error: %s", f.cfg.StreamName, f.cfg.WorkerName, err.Error())
					return
				}
				err := f.seekFile()
				if err !=nil {
					println(err.Error())
				}
				continue
			}
			// 还没有发生mv，说明文件还有可能继续写，这里sleep直到可以读到新的数据
			if same {
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}
	f.dumpRuntime()
}

// 程序起停或者读到已mv文件的文件尾时，寻找下一个要读的文件
// 这里包含核心逻辑:
//			无runtime起停：第一次起
//			有runtime起停：
//				重启时	
//				读到EOF的当前文件已经在backup文件夹中
func (f *FileProducer) seekFile() error {

	// 装填file，更新runtime
	reload := func(path string, offset int64) error {
		file, err := os.Open(path)
		if err != nil {
			log.Error("open file error: %s", err.Error())
			return err
		}
		if _, err := file.Seek(offset, 0); err != nil {
			log.Error("seek error: %s", err.Error())
			return err
		}
		f.offset = offset
		f.reader = util.NewUnitReader(file, f.cfg.Delimiter, f.cfg.BufSize)
		f.runtime = newRuntimeData(file)
		return nil
	}

	//////////////////////////////////////////////////
	//  程序第一次启动，没有runtime数据的情况下     //
	//  PRESENT打开working目录下的文件              //
	//	ORIGINAL从backup目录中找到最早的文件开始读  //
	//////////////////////////////////////////////////
	if f.runtime == nil {
		// 从当前日志文件的文件首开始读
		if f.cfg.Position == PRESENT {
			return reload(f.cfg.FilePath, 0)
		}
		// 找到最早的数据文件，开始读
		return errors.New("ORIGINAL is not support yield")
	}

	////////////////////////////////////////////
	// 第一次启动的情况，reader还没进行初始化 //
	////////////////////////////////////////////
	if f.reader == nil {
		// 比较backup和working中的文件，找到当前文件，确定下一个要读的文件
		if isMatch, _ := f.match(f.cfg.FilePath); isMatch {
			return reload(f.cfg.FilePath, f.runtime.Offset)
		}
		matchFile, err := f.seekFileByInode()
		if err != nil {
			return err
		}
		reload(matchFile, f.runtime.Offset)
		return nil
	}

	/////////////////////////////
	// 需要进行file roll的情况
	// 这部分逻辑简单实现的测试版本
	// 后面需要针对不同的切割逻辑实现
	/////////////////////////////
	matchFile, err := f.seekFileByInode()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(filepath.Base(f.cfg.FilePath) + `\.(\d+)$`)
	ret := re.FindStringSubmatch(matchFile)
	if ret == nil {
		return errors.New("regexp error")
	}
	// 如果存在下一个文件，打开下一个文件进行切割
	postfix, _ := strconv.Atoi(ret[1])
	nextFile := fmt.Sprintf("%s/%s.%d", f.cfg.BackDir, filepath.Base(f.cfg.FilePath), postfix + 1)
	_, err = os.Stat(nextFile)
	if err == nil {
		return reload(nextFile, 0)
	}

	for f.IsActive() && reload(f.cfg.FilePath, 0) != nil {
		time.Sleep(3 * time.Second)
		continue
	}
	return nil
}

// 根据当前的runtime信息，从backup中找到对应的文件的文件名
func (f FileProducer) seekFileByInode() (string, error) {
	pattern := fmt.Sprintf("%s/%s*", f.cfg.BackDir, filepath.Base(f.cfg.FilePath))
	matchFiles, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	for _, matchFile := range matchFiles {
		if isMatch, _ := f.match(matchFile); isMatch {
			return matchFile, nil
		}
	}
	return "", errors.New("can not find the file in backup")
}


// 校验给定路径的文件是不是producer正在读的文件
func (f FileProducer) match(filePath string) (bool, error) {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return fileStat.Sys().(*syscall.Stat_t).Dev == f.runtime.Dev &&
		fileStat.Sys().(*syscall.Stat_t).Ino == f.runtime.Ino, nil
}

type runtimeData struct {
	Offset      int64
	Dev			uint64
	Ino         uint64
}

// producer的runtime数据，这里包含对应的文件信息，并不包含offset信息
// offset只在dump时取即可
func newRuntimeData(file *os.File) *runtimeData {
	// newRuntimeData的情况下，一定是已经打开了文件，所以这里已定不会有异常
	stat, _ := file.Stat()
	data := &runtimeData{
		Dev: stat.Sys().(*syscall.Stat_t).Dev,
		Ino: stat.Sys().(*syscall.Stat_t).Ino,
	}
	return data
}

// dump producer的runtime到特定的文件中
func (f FileProducer) dumpRuntime() {
	runtimeFilePath := fmt.Sprintf("runtime/stream_%s_%s.data", f.cfg.StreamName, f.cfg.WorkerName)
	// 如果文件已经存在，Create函数会覆盖老的文件
	runtimeFile, err := os.Create(runtimeFilePath)
	if err != nil {
		log.Error("%s-%s create data file %s error:", f.cfg.StreamName, f.cfg.WorkerName, runtimeFilePath, err.Error())
		return
	}
	// 这里需要获取当前的offset
	f.runtime.Offset = f.offset
	dataJson, _ := json.Marshal(f.runtime)
	runtimeFile.Write(dataJson)
}

// 程序启动时加载运行时文件，找到上次读的文件
// 在runtimeFile存在的情况下，遇到任何异常都会导致producer初始化失败
func (f *FileProducer) loadRuntime() error {
	runtimeFilePath := fmt.Sprintf("runtime/stream_%s_%s.data", f.cfg.StreamName, f.cfg.WorkerName)
	runtimeFile, err := os.Open(runtimeFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// ReadAll it does not treat an EOF from Read as an error to be reported.
	buf, err := ioutil.ReadAll(runtimeFile)
	if err != nil {
		return err
	}
	var runtime runtimeData
	if err := json.Unmarshal(buf, &runtime); err != nil {
		return err
	}
	f.runtime = &runtime
	return nil
}
