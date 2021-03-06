package sched

import (
	"context"
	"cycron/mod"
	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Job struct {
	taskId       int64
	taskName     string
	expr         *cronexpr.Expression // 解析好的cronexpr表达式
	concurrent   int                  // 并发执行控制
	command      string               // 需要调度的命令
	nextTime     time.Time            // 下次调度时间
	cancelCtx    context.Context      // 任务command的context
	cancelFunc   context.CancelFunc   // 用于取消command执行的cancel函数
	runningCount int                  // 当前正在执行的作业数
	notify       int                  // 是否发送邮件通知
	notifyEmail  []string             // 邮件通知列表
	timeout      int                  // 任务执行超时设置
}

func NewJob(task *mod.TaskMod) (job *Job, err error) {
	var (
		expr       *cronexpr.Expression
		now        time.Time
		ctx        context.Context
		cancelFunc context.CancelFunc
	)

	// 解析JOB的cron表达式
	if expr, err = cronexpr.Parse(task.CronSpec); err != nil {
		log.Errorln(task.TaskName, "表达式解析有误：", err)
		return nil, err
	}

	now = time.Now()

	// 强杀任务时需要用到
	ctx, cancelFunc = context.WithCancel(context.TODO())

	job = &Job{
		taskId:       task.Id,
		taskName:     task.TaskName,
		expr:         expr,
		concurrent:   task.Concurrent,
		command:      task.Command,
		nextTime:     expr.Next(now),
		cancelCtx:    ctx,
		cancelFunc:   cancelFunc,
		runningCount: 0,
		notify:       task.Notify,
		notifyEmail:  nil,
		timeout:      task.Timeout,
	}

	// 至少得允许一个协程运行
	if job.concurrent < 1 {
		job.concurrent = 1
	}

	job.notifyEmail = strings.Split(task.NotifyEmail, "\n")

	return
}

func InitFromTasks(tasks []*mod.TaskMod) (jobs map[int64]*Job, err error) {
	var (
		task *mod.TaskMod
		job  *Job
	)

	if len(tasks) == 0 {
		return
	}

	jobs = make(map[int64]*Job)

	for _, task = range tasks {
		if job, err = NewJob(task); err != nil {
			continue
		}

		jobs[task.Id] = job
	}

	return
}
