package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/huhr/simplelog"
)

const VERSION="0.1.0"

var (
	flagExeDir string
	flagVersion bool
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	initLog()
	parseFlag()
	if flagVersion {
		fmt.Printf("Version: %s\n", VERSION)
		return
	}
	if flagExeDir == "" {
		log.Error("ExeDir is empty")
		exitCode = 1
		return
	}
	if _, err := os.Stat(flagExeDir); err != nil && os.IsNotExist(err) {
		log.Error("ExeDir is illegal")
		exitCode = 1
		return
	}

	os.Chdir(flagExeDir)
	// runtime文件夹用来存放中间文件，中间文件一定要存文本格式，便于特殊需求下手动修改
	if stat, err := os.Stat("runtime"); os.IsNotExist(err) || !stat.IsDir() {
		if err := os.Mkdir("runtime", os.ModePerm); err != nil {
			log.Error(err.Error())
			exitCode = 1
			return
		}
	}
	adapter := NewAdapter()
	exitCode = adapter.Run()
}

func parseFlag() {
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.StringVar(&flagExeDir, "d", "", "Config File")
	flag.Parse()
}

func initLog() {
	log.LoadConfigMap(
		map[string][]map[string]string{
			"root": []map[string]string{
				map[string]string{
					"Level": "info, debug, warn",
					"Output": "log/adapter.log",
					"Rotation":"daily",
					"Format": "detail",
				},
				map[string]string{
					"Level": "error, fatal",
					"Output": "log/adapter.err",
					"Rotation":"daily",
					"Format": "detail",
				},
			},
	})
}

