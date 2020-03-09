package main

import (
	"cycron/sched"
	"fmt"
)

func main()  {
	fmt.Println("这是一个定时任务管理程序")

	sched.Init()
}