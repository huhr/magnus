package main

import (
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/config"
	"github.com/huhr/magnus/stream"
)

// 调度中心
type Adapter struct {
	exeDir string
	streams []*stream.Stream
}

func NewAdapter(exeDir string) *Adapter {
	return &Adapter{exeDir: exeDir}
}

// 加载配置文件，创建各个stream
func (a *Adapter) initStream() error {
	files, _ := filepath.Glob(a.exeDir+"/conf/stream_*.toml")
	for _, file := range files {
		var cfg config.StreamConfig
		if _, err := toml.DecodeFile(file, &cfg); err != nil {
			log.Error(err.Error())
			return nil
		}
		a.streams = append(a.streams, stream.NewStream(cfg))
	}
	return nil
}


// 启动各个stream
func (a Adapter) Run() int {
	var wg sync.WaitGroup

	log.Debug("init streams")
	if a.initStream() != nil {
		return 1
	}
	log.Debug("total %d streams", len(a.streams))
	for _, s := range a.streams {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Run()
		}()
	}
	wg.Wait()
	return 0
}

