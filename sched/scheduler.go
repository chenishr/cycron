package sched

import (
	"cycron/mod"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Scheduler struct {
	jobs         map[int64]*Job   // 需要调度的作业集合
	addJobChan   chan bool        // 通知调度器有新任务添加
	resChan      chan *ExecResult // 作业执行结果
	running      bool             // 调度器是否已经启动
	runningCount int64            // 当前正在执行的任务数
}

// 任务执行结果
type ExecResult struct {
	job       *Job      // 执行的任务
	output    []byte    // 脚本输出
	err       error     // 脚本错误原因
	status    int       //
	planTime  time.Time // 理论上的调度时间
	realTime  time.Time // 实际的调度时间
	startTime time.Time // 启动时间
	endTime   time.Time // 结束时间
}

var (
	GScheduler *Scheduler
)

func InitScheduler() {
	var (
		s *Scheduler
	)

	s = &Scheduler{
		jobs:    nil,
		resChan: nil,
		running: false,
	}

	s.resChan = make(chan *ExecResult, 100)
	s.jobs = make(map[int64]*Job, 100)
	s.addJobChan = make(chan bool, 10)
	log.Info("初始化调度队列")

	GScheduler = s
}

func (s *Scheduler) StartScheduler() (err error) {
	var (
		tasks []*mod.TaskMod
		jobs  map[int64]*Job
	)

	// 获取需要调度的任务
	tasks, err = mod.GTaskMgr.GetTasks()
	if err != nil {
		return
	}

	// 初始化调度队列
	jobs, err = InitFromTasks(tasks)
	if err != nil {
		return
	}

	// 任务不为空才进行赋值
	if jobs != nil {
		s.jobs = jobs
	}

	// 启动任务调度循环
	s.loop()

	return
}

func (s *Scheduler) RemoveJob(taskId int64) {
	var (
		exists bool
		oldJob *Job
	)
	if s.running == false {
		return
	}

	// 检查该任务是否在调度队列中
	if oldJob, exists = s.jobs[taskId]; !exists {
		return
	}

	// 终止正在调度的任务
	if oldJob.runningCount > 0 {
		log.Warnln("触发取消", oldJob.taskName, "执行中的", oldJob.runningCount, "个任务")
		oldJob.cancelFunc()
	}

	// 从调度队列中删除任务
	delete(s.jobs, taskId)

	return
}

func (s *Scheduler) AddJob(task *mod.TaskMod, readd bool) (err error) {
	var (
		exists bool
		job    *Job
	)
	if s.running == false {
		return
	}

	// 检查该任务是否在调度队列中
	if _, exists = s.jobs[task.Id]; exists == true {
		if readd {
			// 是否需要丢弃从新加入
			s.RemoveJob(task.Id)
		} else {
			return
		}
	}

	if job, err = NewJob(task); err != nil {
		return
	}

	s.jobs[task.Id] = job

	s.addJobChan <- true

	return
}

func (s *Scheduler) loop() {
	if s.running {
		return
	}

	s.running = true

	// 启动事件处理协程
	s.HandleEvent()

	var (
		now      time.Time
		nearTime *time.Time
		waitTime time.Duration
	)

	for {
		now = time.Now()

		log.Traceln("当前调度时间：", now)

		if len(s.jobs) == 0 {
			log.Info("暂时没有任务执行", now)
			waitTime = 1 * time.Second
		}

		for _, job := range s.jobs {
			if job.nextTime.Before(now) || job.nextTime.Equal(now) {

				res := &ExecResult{
					job:       job,
					output:    nil,
					err:       nil,
					status:    mod.TASK_SUCCESS,
					planTime:  job.nextTime,
					realTime:  now,
					startTime: time.Time{},
					endTime:   time.Time{},
				}

				// 控制并发
				if job.runningCount < job.concurrent {
					job.runningCount++
					s.runningCount++

					log.Traceln(job.taskName, "当前进行作业数", job.runningCount)
					log.Traceln("调度器当前进行作业数", s.runningCount)

					// 执行任务
					GExecutor.ExecuteJob(res)
				} else {
					log.Warnln("任务[", job.taskName, "]协程启动数量将超过允许的", job.concurrent, "个，本次被忽略")
					res.status = mod.TASK_IGNORE
					s.PushJobResult(res)
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
			log.Info("下次调度时间：" ,nearTime)
			log.Info()

		*/

		// 睡眠100毫秒
		select {
		case <-time.NewTimer(waitTime).C: // 将在100毫秒可读，返回
			nearTime = nil
		case <-s.addJobChan:
			// 有新任务到来，重新计算调度时间
		}
	}
}

func (s *Scheduler) RunOnce(taskId int64) {
	var (
		job    *Job
		exists bool
	)
	if s.running == false {
		return
	}

	if job, exists = s.jobs[taskId]; !exists {
		return
	}

	s.RunJob(job)
}

func (s *Scheduler) RunJob(job *Job) {
	// 控制并发
	job.runningCount++
	s.runningCount++

	log.Traceln(job.taskName, "当前进行作业数", job.runningCount)

	// 执行任务
	res := &ExecResult{
		job:       job,
		output:    nil,
		err:       nil,
		status:    mod.TASK_SUCCESS,
		planTime:  job.nextTime,
		realTime:  time.Now(),
		startTime: time.Time{},
		endTime:   time.Time{},
	}
	GExecutor.ExecuteJob(res)
}

func (s *Scheduler) HandleEvent() {
	go func() {
		var (
			errMsg  string
			uptCond bson.M
			uptData bson.M
		)

		for {
			// 睡眠100毫秒
			select {
			case res := <-s.resChan:
				res.job.runningCount--
				s.runningCount--

				if res.err != nil {
					errMsg = res.err.Error()
				} else {
					errMsg = ""
				}
				log.Debugln(res.job.taskName + "打印结果:" + string(res.output) + "打印错误：" + errMsg)

				// 发邮件通知
				if (res.job.notify == 1 && res.err != nil) || res.job.notify == 2 {
					GMailer.OrgData(res)
				}

				// 保存执行日记
				GLogger.OrgData(res)

				// 保存上次执行时间
				uptCond = bson.M{"_id": res.job.taskId}
				uptData = bson.M{
					"$set": bson.M{
						"prev_time": res.endTime.Format("2006-01-02 15:04:05"),
					},
				}
				mod.GTaskMgr.UpdateOne(uptCond, uptData)
			}
		}
	}()
}

// 回传任务执行结果
func (s *Scheduler) PushJobResult(execResult *ExecResult) {
	s.resChan <- execResult
}
