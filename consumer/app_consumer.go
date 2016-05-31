package consumer

import (
	"io"
	"os"
	"os/exec"

	log "github.com/huhr/simplelog"

	"github.com/huhr/magnus/filter"
)

// 启动下游程序，并将数据发送到下游程序的stdin
type AppConsumer struct {
	BaseConsumer
	cmd			*exec.Cmd
	appStdin	io.WriteCloser
}

func NewAppConsumer(base BaseConsumer) *AppConsumer {
	cmd := exec.Command(base.cfg.StartupScript)
	appStdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}
	// 重定向标准输出和标准错误输出
	cmd.Stdout, err = os.Create(base.cfg.OutputDir + "/stdout")
	if err != nil {
		return nil
	}
	cmd.Stderr, err = os.Create(base.cfg.OutputDir + "/stderr")
	if err != nil {
		return nil
	}
	err = cmd.Start()
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &AppConsumer{base, cmd, appStdin}
}

func (app *AppConsumer) Consume(msg []byte) bool {
	if !filter.Filter(msg, app.cfg.Filters) {
		if _, err := app.appStdin.Write(msg); err != nil{
			log.Error(err.Error())
			return false
		}
	}
	return true
}
