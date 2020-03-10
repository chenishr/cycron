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
	Output 		[]byte 		// 脚本输出
	Err 		error 		// 脚本错误原因
	StartTime 	time.Time 	// 启动时间
	EndTime 	time.Time 	// 结束时间
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

	tasks,err = models.GetTasks()
	if err != nil{
		return
	}

	s.jobs,err = InitFromTasks(tasks)

	return
}
