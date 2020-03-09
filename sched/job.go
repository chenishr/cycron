package sched

import (
	"context"
	"fmt"
	"time"
	"github.com/gorhill/cronexpr"
)

type Job struct {
	taskId 		int
	taskName	string
	expr 		*cronexpr.Expression	// 解析好的cronexpr表达式
	concurrent  int						// 并发执行控制
	command     string					// 需要调度的命令
	nextTime 	time.Time				// 下次调度时间
	cancelCtx 	context.Context 		// 任务command的context
	cancelFunc 	context.CancelFunc		// 用于取消command执行的cancel函数
}

func Init()  {
	fmt.Println("hells")
}