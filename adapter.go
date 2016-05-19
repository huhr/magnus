package main

import (
	"path/filepath"
	"os"
	"os/signal"
	"syscall"
	"sync"

	"github.com/BurntSushi/toml"
	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/config"
	"github.com/huhr/magnus/stream"
)

// 调度中心
type Adapter struct {
	streams []*stream.Stream
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

// 加载配置文件，创建各个stream
func (a *Adapter) initStream() error {
	files, err := filepath.Glob("conf/stream_*.toml")
	if err != nil {
		log.Error("Glob conf files error: %s", err.Error())
		return err
	}
	for _, file := range files {
		var cfg config.StreamConfig
		if _, err := toml.DecodeFile(file, &cfg); err != nil {
			log.Error("Decode toml file error: %s", err.Error())
			return err
		}
		a.streams = append(a.streams, stream.NewStream(cfg))
	}
	return nil
}


// 启动各个stream
func (a Adapter) Run() int {
	var wg sync.WaitGroup

	a.registerSigalHandler()
	log.Debug("begin to init streams")
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

func (a *Adapter) registerSigalHandler() {
    go func() {
        for {
            c := make(chan os.Signal)
            signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
            // sig is blocked as c is 没缓冲
            sig := <-c
            log.Info("Signal %d received", sig)
			for _, s := range a.streams {
				s.ShutDown()
			}
        }
    }()
}

