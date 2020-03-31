package main

import (
	"cycron/api"
	"cycron/sched"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	var (
		err error
	)
	fmt.Println("这是一个定时任务管理程序")

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	log.SetLevel(log.TraceLevel)
	log.Info("slkdjfsd")
	// 启动 HTTPServer
	err = api.InitHttpServer()
	if err != nil {
		fmt.Println("HttpServer启动失败：", err)
	}

	err = sched.GScheduler.InitScheduler()
	if err != nil {
		fmt.Println("调度器启动失败：", err)
	}

	sched.GScheduler.Print()

	time.Sleep(3600 * time.Second)
	fmt.Println("一小时后程序退出")
}
