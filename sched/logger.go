package sched

import (
	"cycron/conf"
	"cycron/mod"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type Logger struct {
	logChan        chan *mod.TaskLogMod
	autoCommitChan chan *LogBatch
	statChan       chan *mod.TaskLogStatMod
	statCommitChan chan map[string]*mod.TaskLogStatMod
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
		log.Errorln("保存日志失败：", err)
	}
}

// 日志存储协程
func (l *Logger) statLoop() {
	var (
		statData    *mod.TaskLogStatMod
		statMap     map[string]*mod.TaskLogStatMod
		timeoutStat map[string]*mod.TaskLogStatMod
		key         string
		ok          bool
		mux         sync.Mutex
	)

	for {
		select {
		case statData = <-l.statChan:
			if statMap == nil {
				statMap = make(map[string]*mod.TaskLogStatMod)
				// 让这个批次超时自动提交(给1秒的时间）
				time.AfterFunc(
					time.Duration(conf.GConfig.Mongo.CommitTimeout*5)*time.Millisecond,
					func(stat map[string]*mod.TaskLogStatMod) func() {
						return func() {
							l.statCommitChan <- stat
						}
					}(statMap),
				)
			}

			// 统计
			key = statData.PlanTime + "::" + strconv.FormatInt(int64(statData.Status), 10)
			if _, ok = statMap[key]; ok {
				statMap[key].Count++
			} else {
				statMap[key] = &mod.TaskLogStatMod{
					Status:   statData.Status,
					PlanTime: statData.PlanTime,
					Count:    1,
				}
			}

		case timeoutStat = <-l.statCommitChan: // 过期的批次
			mux.Lock()
			// 把批次写入到mongo中
			for _, val := range timeoutStat {
				mod.GTaskLogStatMgr.UpsertDoc(val)
			}
			// 清空logBatch
			statMap = nil

			mux.Unlock()
		}
	}
}

// 日志存储协程
func (l *Logger) writeLoop() {
	var (
		logData      *mod.TaskLogMod
		logBatch     *LogBatch // 当前的批次
		commitTimer  *time.Timer
		timeoutBatch *LogBatch // 超时批次
	)

	for {
		select {
		case logData = <-l.logChan:
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
			logBatch.Logs = append(logBatch.Logs, logData)

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
		statChan:       make(chan *mod.TaskLogStatMod, 1000),
		statCommitChan: make(chan map[string]*mod.TaskLogStatMod, 1000),
	}

	// 启动一个mongodb处理协程
	go GLogger.writeLoop()
	go GLogger.statLoop()
	return
}

// 发送日志
func (l *Logger) Append(taskLog *mod.TaskLogMod) {
	select {
	case l.logChan <- taskLog:
	default:
		// 队列满了就丢弃
		log.Warnln("日记队列已满，丢弃本次执行日记")
	}

	t, _ := time.Parse("2006-01-02 15:04:05", taskLog.PlanTime)
	// 记录精度，小时
	plantime := t.Format("2006-01-02 15")
	stat := &mod.TaskLogStatMod{
		Status:   taskLog.Status,
		PlanTime: plantime,
		Count:    1,
	}
	select {
	case l.statChan <- stat:
	default:
		// 队列满了就丢弃
		log.Warnln("日记统计队列已满，丢弃本次执行日记")
	}
}

func (l *Logger) OrgData(res *ExecResult) {
	var (
		taskLog *mod.TaskLogMod
		errMsg  string
		err     error
		id      int64
		psTime  int
	)

	if mod.TASK_IGNORE != res.status {
		psTime = int(res.endTime.Sub(res.realTime) / time.Millisecond)
		if res.err != nil {
			errMsg = res.err.Error()
		} else {
			errMsg = ""
		}
	}

	id, err = mod.GCommonMgr.GetMaxId(conf.GConfig.Models.TaskLog)
	if err != nil {
		log.Errorln("获取"+conf.GConfig.Models.TaskLog+"自增 ID 失败", err)
		return
	}
	taskLog = &mod.TaskLogMod{
		Id:          id,
		TaskId:      res.job.taskId,
		Command:     res.job.command,
		Status:      res.status,
		ProcessTime: psTime,
		Output:      string(res.output),
		Err:         errMsg,
		PlanTime:    res.planTime.Format("2006-01-02 15:04:05"),
		RealTime:    res.realTime.Format("2006-01-02 15:04:05"),
		StartTime:   res.startTime.Format("2006-01-02 15:04:05"),
		EndTime:     res.endTime.Format("2006-01-02 15:04:05"),
	}

	l.Append(taskLog)
}
