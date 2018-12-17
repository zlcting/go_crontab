package main

import (
	"flag"
	"fmt"
	"runtime"
	"zlc_sys/master"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var confFile string

func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "传入master.json")
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
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}
	//启动APi http 服务

	if err = master.InitApiServer(); err != nil {
		goto ERR

	}

	return
ERR:
	fmt.Println(err)
}
