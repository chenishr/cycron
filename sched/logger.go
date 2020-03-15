package sched

import (
	"cycron/conf"
	"cycron/mod"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Logger struct {
	logChan        chan *mod.TaskLogMod
	autoCommitChan chan *LogBatch
}

// 日志批次
type LogBatch struct {
	Logs []interface{} // 多条日志
}

var (
	// 单例
	GLogger *Logger
)

// 批量写入日志
func (l *Logger) saveLogs(batch *LogBatch) {
	err := mod.GTaskLogMgr.InsertMany(batch.Logs)
	if err != nil {
		fmt.Println("保存日志失败：", err)
	}
}

// 日志存储协程
func (l *Logger) writeLoop() {
	var (
		log          *mod.TaskLogMod
		logBatch     *LogBatch // 当前的批次
		commitTimer  *time.Timer
		timeoutBatch *LogBatch // 超时批次
	)

	for {
		select {
		case log = <-l.logChan:
			if logBatch == nil {
				logBatch = &LogBatch{}
				// 让这个批次超时自动提交(给1秒的时间）
				commitTimer = time.AfterFunc(
					time.Duration(conf.GConfig.Mongo.CommitTimeout)*time.Millisecond,
					func(batch *LogBatch) func() {
						return func() {
							l.autoCommitChan <- batch
						}
					}(logBatch),
				)
			}

			// 把新日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果批次满了, 就立即发送
			if len(logBatch.Logs) >= conf.GConfig.Logger.BatchSize {
				// 发送日志
				l.saveLogs(logBatch)
				// 清空logBatch
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <-l.autoCommitChan: // 过期的批次
			// 判断过期批次是否仍旧是当前的批次
			if timeoutBatch != logBatch {
				continue // 跳过已经被提交的批次
			}
			// 把批次写入到mongo中
			l.saveLogs(timeoutBatch)
			// 清空logBatch
			logBatch = nil
		}
	}
}

func init() {
	GLogger = &Logger{
		logChan:        make(chan *mod.TaskLogMod, 1000),
		autoCommitChan: make(chan *LogBatch, 1000),
	}

	// 启动一个mongodb处理协程
	go GLogger.writeLoop()
	return
}

// 发送日志
func (l *Logger) Append(taskLog *mod.TaskLogMod) {
	select {
	case l.logChan <- taskLog:
	default:
		// 队列满了就丢弃
		fmt.Println("日记队列已满，丢弃本次执行日记")
	}
}

func (l *Logger) OrgData(res *ExecResult) {
	var (
		taskLog *mod.TaskLogMod
		errMsg  string
	)

	psTime := int(res.endTime.Sub(res.realTime) / time.Millisecond)
	if res.err != nil {
		errMsg = res.err.Error()
	} else {
		errMsg = ""
	}

	id, _ := primitive.ObjectIDFromHex(res.job.taskId)
	taskLog = &mod.TaskLogMod{
		Id:          primitive.NewObjectID(),
		TaskId:      id,
		Command:     res.job.command,
		Status:      res.status,
		ProcessTime: psTime,
		Output:      string(res.output),
		Err:         errMsg,
		PlanTime:    res.planTime.Unix(),
		RealTime:    res.realTime.Unix(),
		StartTime:   res.startTime.Unix(),
		EndTime:     res.endTime.Unix(),
	}

	l.Append(taskLog)
}
