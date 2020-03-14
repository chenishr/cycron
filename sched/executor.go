package sched

import (
	"cycron/models"
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
func (executor *Executor) ExecuteJob(result *ExecResult) {
	go func() {
		var (
			cmd *exec.Cmd
			err error
			output []byte
		)

		result.startTime 	= time.Now()
		result.endTime		= time.Now()

		// 任务结果
		// 执行shell命令
		fmt.Println("执行器开始执行" ,result.job.taskName,time.Now())
		cmd = exec.CommandContext(result.job.cancelCtx,"/bin/bash","-c",result.job.command)

		// 执行并捕获输出
		output, err = cmd.CombinedOutput()
		fmt.Println("执行器结束执行" ,result.job.taskName,time.Now())

		// 记录任务结束时间
		result.endTime = time.Now()
		result.output = output
		result.err = err

		if err != nil{
			result.status = models.TASK_ERROR
		}

		// 上报执行结果
		GScheduler.PushJobResult(result)
	}()
}

//  初始化执行器
func InitExecutor() (err error) {
	GExecutor = &Executor{}
	return
}
