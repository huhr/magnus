package main

import (
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/BurntSushi/toml"
	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/stream"
	"github.com/huhr/magnus/tools"
)

// Adapter负责调度各个stream进行运转
type Adapter struct {
	streams []*stream.Stream
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) initStream() error {
	log.Debug("Degin To Init Streams")
	files, err := filepath.Glob("conf/stream_*.toml")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("Can Not Find stream_*.toml File")
	}

	var config tools.StreamConfig
	for _, file := range files {
		if _, err := toml.DecodeFile(file, &config); err != nil {
			log.Error("Decode Toml File Error: %s", err.Error())
			continue
		}
		s, err := stream.NewStream(config)
		if err != nil {
			log.Error("%s Init Error: %s", config.StreamName, err.Error())
			continue
		}
		a.streams = append(a.streams, s)
	}

	if len(a.streams) == 0 {
		return errors.New("Init No Success Streams")
	}
	log.Debug("Init %d Streams", len(a.streams))
	return nil
}

// 注册监听Control-C、SIGTERM信号，优雅停止
// 监听其他信号，处理如reload等操作
func (a *Adapter) registerSigalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Info("Signal %d received", sig)
	for _, s := range a.streams {
		s.ShutDown()
	}
}

func (a Adapter) Run() int {
	go a.registerSigalHandler()
	if err := a.initStream(); err != nil {
		log.Error(err.Error())
		return 1
	}

	var wg sync.WaitGroup
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
