package sched

import (
	"context"
	"cycron/conf"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Logger struct {
	client 			*mongo.Client
	logCollection 	*mongo.Collection
	logChan 		chan *TaskLog
	autoCommitChan 	chan *LogBatch
}

// 任务执行日志
type TaskLog struct {
	TaskId      	int		`json:"task_id" bson:"task_id"`
	Command 		string 	`json:"command" bson:"command"` // 脚本命令
	Status      	int		`json:"status" bson:"status"`
	ProcessTime 	int		`json:"process_time" bson:"process_time"`
	Output 			string 	`json:"output" bson:"output"`	// 脚本输出
	Err 			string 	`json:"err" bson:"err"` // 错误原因
	PlanTime 		int64 	`json:"plan_time" bson:"plan_time"`	// 理论上的调度时间
	RealTime 		int64 	`json:"real_time" bson:"real_time"`	// 实际的调度时间
	StartTime 		int64 	`json:"start_time" bson:"start_time"` // 启动时间
	EndTime 		int64 	`json:"end_time" bson:"end_time"` 	// 结束时间
}

// 日志批次
type LogBatch struct {
	Logs []interface{}	// 多条日志
}

var (
	// 单例
	GLogger *Logger
)

// 批量写入日志
func (l *Logger) saveLogs(batch *LogBatch) {
	l.logCollection.InsertMany(context.TODO(), batch.Logs)
}

// 日志存储协程
func (l *Logger) writeLoop() {
	var (
		log *TaskLog
		logBatch *LogBatch // 当前的批次
		commitTimer *time.Timer
		timeoutBatch *LogBatch // 超时批次
	)

	for {
		select {
		case log = <- l.logChan:
			if logBatch == nil {
				logBatch = &LogBatch{}
				// 让这个批次超时自动提交(给1秒的时间）
				commitTimer = time.AfterFunc(
					time.Duration(conf.GConfig.Mongo.CommitTimeout) * time.Millisecond,
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
		case timeoutBatch = <- l.autoCommitChan: // 过期的批次
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
	var (
		client *mongo.Client
		mongoConf conf.MongoConf
		loggerConf conf.LoggerConf
		ctx context.Context
		err error
	)

	mongoConf = conf.GConfig.Mongo
	loggerConf = conf.GConfig.Logger

	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(mongoConf.ConnectTimeout) * time.Millisecond)

	// 建立mongodb连接
	fmt.Println("建立mongodb连接")
	if client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoConf.Uri)); err != nil {
		fmt.Println("链接 MongoDB 失败：",err)
		return
	}
	fmt.Println("建立mongodb连接 成功")

	//   选择db和collection
	GLogger = &Logger{
		client: client,
		logCollection: client.Database(loggerConf.Db).Collection(loggerConf.TaskLog),
		logChan: make(chan *TaskLog, 1000),
		autoCommitChan: make(chan *LogBatch, 1000),
	}

	// 启动一个mongodb处理协程
	go GLogger.writeLoop()
	return
}

// 发送日志
func (l *Logger) Append(taskLog *TaskLog) {
	select {
	case l.logChan <- taskLog:
	default:
		// 队列满了就丢弃
		fmt.Println("日记队列已满，丢弃本次执行日记")
	}
}

func (l *Logger) OrgData(res *ExecResult)  {
	var (
		taskLog *TaskLog
		errMsg string
	)

	psTime 	:= int(res.endTime.Sub(res.realTime) / time.Millisecond)
	if res.err != nil {
		errMsg = res.err.Error()
	}else{
		errMsg = ""
	}

	taskLog = &TaskLog{
		TaskId:      res.job.taskId,
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
