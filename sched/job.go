package sched

import (
	"context"
	"cycron/models"
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
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

func InitFromTasks(tasks []*models.TaskMod) (jobs map[int]*Job, err error) {
	var (
		task 	*models.TaskMod
		job		*Job
		expr 	*cronexpr.Expression
		now		time.Time
		ctx		context.Context
		cancelFunc	context.CancelFunc
	)

	if len(tasks) == 0 {
		return
	}

	jobs = make(map[int]*Job)

	now = time.Now()

	for task,_ = range tasks {
		// 解析JOB的cron表达式
		if expr, err = cronexpr.Parse(task.CronSpec); err != nil {
			fmt.Println(task.TaskName,"表达式解析有误：",err)
			continue;
		}

		ctx,cancelFunc = context.WithCancel(context.TODO())

		job = &Job{
			taskId:     task.Id,
			taskName:   task.TaskName,
			expr:       expr,
			concurrent: task.Concurrent,
			command:    task.Command,
			nextTime:   expr.Next(now),
			cancelCtx:  ctx,
			cancelFunc: cancelFunc,
		}
		jobs[task.Id] = job
	}

	return
}