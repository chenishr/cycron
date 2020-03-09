package sched

import "time"

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
