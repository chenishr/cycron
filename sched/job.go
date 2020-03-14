package sched

import (
	"context"
	"cycron/models"
	"fmt"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	taskId 		string
	taskName	string
	expr 		*cronexpr.Expression	// 解析好的cronexpr表达式
	concurrent  int						// 并发执行控制
	command     string					// 需要调度的命令
	nextTime 	time.Time				// 下次调度时间
	cancelCtx 	context.Context 		// 任务command的context
	cancelFunc 	context.CancelFunc		// 用于取消command执行的cancel函数
	runningCount	int					// 当前正在执行的作业数
	notify      	int					// 是否发送邮件通知
	notifyEmail  	[]string			// 邮件通知列表
}

func InitFromTasks(tasks []*models.TaskMod) (jobs map[string]*Job, err error) {
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

	jobs = make(map[string]*Job)

	now = time.Now()

	for _,task = range tasks {
		// 解析JOB的cron表达式
		if expr, err = cronexpr.Parse(task.CronSpec); err != nil {
			fmt.Println(task.TaskName,"表达式解析有误：",err)
			continue;
		}

		// 强杀任务时需要用到
		ctx,cancelFunc = context.WithCancel(context.TODO())

		job = &Job{
			taskId:     task.Id.Hex(),
			taskName:   task.TaskName,
			expr:       expr,
			concurrent: task.Concurrent,
			command:    task.Command,
			nextTime:   expr.Next(now),
			cancelCtx:  ctx,
			cancelFunc: cancelFunc,
			runningCount: 	0,
			notify:			task.Notify,
			notifyEmail:  	nil,
		}

		// 至少得允许一个协程运行
		if job.concurrent < 1 {
			job.concurrent = 1
		}

		job.notifyEmail = strings.Split(task.NotifyEmail, "\n")

		jobs[task.Id.String()] = job
	}

	return
}