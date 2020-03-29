package main

import (
	"cycron/api"
	"cycron/sched"
	"fmt"
	"time"
)

func main() {
	var (
		err error
	)
	fmt.Println("这是一个定时任务管理程序")

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
