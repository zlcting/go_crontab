package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
	"zlc_sys/worker"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var confFile string

func initArgs() {
	//woker -config ./worker.josn
	//worker -h

	flag.StringVar(&confFile, "config", "./worker.json", "传入worker.json")
	flag.Parse()
}

func main() {
	var (
		err error
	)
	//加载命令行参数
	initArgs()
	//初始化线程
	initEnv()
	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}
	//初始化任务管理器
	if err = worker.InitJobMgr(); err != nil {
		goto ERR

	}

	for {
		time.Sleep(1 * time.Second)
	}
	return
ERR:
	fmt.Println(err)
}
