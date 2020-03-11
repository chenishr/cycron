package sched

import (
	"fmt"
	"os/exec"
	"time"
)

// 任务执行器
type Executor struct {

}
var (
	GExecutor *Executor
)

// 执行一个任务
func (executor *Executor) ExecuteJob(job *Job) {
	go func() {
		var (
			result *ExecResult
			cmd *exec.Cmd
			err error
			output []byte
		)

		result = &ExecResult{
			job:		job,
			output:    	nil,
			err:       	nil,
			startTime: 	time.Now(),
			endTime:   	time.Now(),
		}

		// 任务结果
		// 执行shell命令
		fmt.Println("执行器开始执行" ,job.taskName,time.Now())
		cmd = exec.CommandContext(job.cancelCtx,"/bin/bash","-c",job.command)

		// 执行并捕获输出
		output, err = cmd.CombinedOutput()
		fmt.Println("执行器结束执行" ,job.taskName,time.Now())

		// 记录任务结束时间
		result.endTime = time.Now()
		result.output = output
		result.err = err

		// 上报执行结果
		GScheduler.PushJobResult(result)
	}()
}

//  初始化执行器
func InitExecutor() (err error) {
	GExecutor = &Executor{}
	return
}
