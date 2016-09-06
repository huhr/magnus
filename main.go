package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/huhr/simplelog"
)

const VERSION = "0.1.0"

var (
	flagExeDir  string
	flagVersion bool
)

func main() {
	exitCode := 0
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintln(os.Stderr, p)
			os.Exit(1)
		}
		os.Exit(exitCode)
	}()

	initLog()

	parseFlag()
	if flagVersion {
		fmt.Printf("Version: %s\n", VERSION)
		return
	}
	if flagExeDir != "" {
		if _, err := os.Stat(flagExeDir); err != nil && os.IsNotExist(err) {
			log.Error("ExeDir is illegal")
			exitCode = 1
			return
		}
		os.Chdir(flagExeDir)
	}
	// runtime文件夹用来存放中间文件，中间文件存文本格式，便于特殊需求下手动修改
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

// 命令行工具
func parseFlag() {
	flag.BoolVar(&flagVersion, "v", false, "Version")
	flag.StringVar(&flagExeDir, "d", "", "ExeDir")
	flag.Parse()
}

// 初始化日志配置，使用了simplelog
func initLog() {
	err := log.LoadConfigMap(
		map[string][]map[string]string{
			"root": []map[string]string{
				map[string]string{
					"Level":    "info, debug, warn",
					"Output":   "logs/adapter.log",
					"Rotation": "daily",
					"Format":   "detail",
				},
				map[string]string{
					"Level":    "error, fatal",
					"Output":   "logs/adapter.err",
					"Rotation": "daily",
					"Format":   "detail",
				},
			},
		})
	if err != nil {
		fmt.Println("asd")
	}
}
