package consumer

import (
	"io"
	"os"
	"os/exec"

	log "github.com/huhr/simplelog"
)

// 启动下游程序，并将数据发送到下游程序的stdin
type AppConsumer struct {
	BaseConsumer
	cmd      *exec.Cmd
	appStdin io.WriteCloser
}

func NewAppConsumer(base BaseConsumer) (Consumer, error) {
	cmd := exec.Command(base.config.StartupScript)
	appStdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	// 重定向标准输出和标准错误输出
	cmd.Stdout, err = os.Create(base.config.OutputDir + "/stdout")
	if err != nil {
		return nil, err
	}
	cmd.Stderr, err = os.Create(base.config.OutputDir + "/stderr")
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return &AppConsumer{base, cmd, appStdin}, nil
}

func (app *AppConsumer) Consume(msg []byte) bool {
	if _, err := app.appStdin.Write(msg); err != nil {
		log.Error(err.Error())
		return false
	}
	return true
}
