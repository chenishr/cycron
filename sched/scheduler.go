package sched

import (
	"cycron/models"
	"fmt"
	"time"
)

type Scheduler struct {
	jobs	map[int]*Job		// 需要调度的作业集合
	resChan	chan *ExecResult	// 作业执行结果
	running	bool				// 调度器是否已经启动
}

// 任务执行结果
type ExecResult struct {
	job		*Job		// 执行的任务
	output 		[]byte 		// 脚本输出
	err 		error 		// 脚本错误原因
	planTime 	time.Time 	// 理论上的调度时间
	realTime 	time.Time 	// 实际的调度时间
	startTime 	time.Time 	// 启动时间
	endTime 	time.Time 	// 结束时间
}

var(
	GScheduler * Scheduler
)

func init()  {
	var (
		s 		*Scheduler
		err		error
	)

	s = &Scheduler{
		jobs:    nil,
		resChan: nil,
		running: false,
	}

	s.resChan = make(chan *ExecResult,100)

	err = s.InitScheduler()
	if err != nil {
		fmt.Println("初始化调度器失败",err)
	}

	GScheduler = s

	// 启动任务调度循环
	s.loop()
}

func (s *Scheduler)Print() {
	for id,job := range s.jobs {
		fmt.Println(id,job)
	}
}

func (s *Scheduler)InitScheduler()(err error) {
	var (
		tasks 	[]*models.TaskMod
	)

	// 获取需要调度的任务
	tasks,err = models.GetTasks()
	if err != nil{
		return
	}

	// 初始化调度队列
	s.jobs,err = InitFromTasks(tasks)
	if err != nil {
		return
	}

	return
}

func (s *Scheduler)loop() {
	if s.running {
		return
	}

	s.running = true

	// 启动事件处理协程
	s.HandleEvent()

	var(
		now	time.Time
		nearTime *time.Time
		waitTime time.Duration
	)

	for{
		now = time.Now()

		//fmt.Println("当前调度时间：" ,now)

		if len(s.jobs) == 0 {
			waitTime = 1 * time.Second
		}

		for _, job := range s.jobs {
			if job.nextTime.Before(now) || job.nextTime.Equal(now) {

				// 控制并发
				if job.runningCount < job.concurrent {
					job.runningCount ++

					fmt.Println(job.taskName, "当前进行作业数", job.runningCount)

					// 执行任务
					res := &ExecResult{
						job:      	job,
						output:    nil,
						err:       nil,
						planTime: job.nextTime,
						realTime: now,
						startTime: time.Time{},
						endTime:   time.Time{},
					}
					GExecutor.ExecuteJob(res)
				} else {
					fmt.Println("任务[",job.taskName,"]协程启动数量将超过允许的", job.concurrent ,"个，本次被忽略")
				}

				// 每次调度之后都需要更新下次执行时间
				job.nextTime = job.expr.Next(now)
			}

			if nearTime == nil || job.nextTime.Before(*nearTime) {
				nearTime = &job.nextTime
				waitTime = nearTime.Sub(now)
			}

		}

		/*
		fmt.Println("下次调度时间：" ,nearTime)
		fmt.Println()

		 */

		// 睡眠100毫秒
		select {
		case <- time.NewTimer(waitTime).C:	// 将在100毫秒可读，返回
			nearTime = nil
		}
	}
}

func (s *Scheduler)HandleEvent()  {
	go func() {
		var(
			errMsg string
		)

		for {
			// 睡眠100毫秒
			select {
			case res := <-s.resChan:
				res.job.runningCount --

				if res.err != nil {
					errMsg = res.err.Error()
				} else {
					errMsg = ""
				}
				fmt.Println(res.job.taskName + "打印结果:" + string(res.output) + "打印错误：" + errMsg)

				// 发邮件通知
				if (res.job.notify == 1 && res.err != nil) || res.job.notify == 2 {
					GMailer.OrgData(res)
				}
			}
		}
	}()
}

// 回传任务执行结果
func (s *Scheduler) PushJobResult(execResult *ExecResult) {
	s.resChan <- execResult
}
