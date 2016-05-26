package consumer

import (
	"fmt"
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
	appStdout   io.ReadCloser
}

func NewAppConsumer(base BaseConsumer) *AppConsumer {
	cmd := exec.Command(base.cfg.StartupScript)
	appStdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}
	appStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}
	err = cmd.Start()
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	go func() {
		io.Copy(os.Stdin, appStdout)
	}()
	return &AppConsumer{base, cmd, appStdin, appStdout}
}

func (app *AppConsumer) Consume(msg []byte) bool {
	if !filter.Filter(msg, app.cfg.Filters) {
		fmt.Println("msg: %s", msg)
		n, err := app.appStdin.Write(msg)
		fmt.Println(n)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return true
}
