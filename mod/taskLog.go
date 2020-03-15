package mod

import "go.mongodb.org/mongo-driver/bson/primitive"

// 任务执行日志
type TaskLogMod struct {
	Id          primitive.ObjectID `json:"task_id" bson:"_id"`
	TaskId      primitive.ObjectID `json:"task_id" bson:"task_id"`
	Command     string             `json:"command" bson:"command"` // 脚本命令
	Status      int                `json:"status" bson:"status"`
	ProcessTime int                `json:"process_time" bson:"process_time"`
	Output      string             `json:"output" bson:"output"`         // 脚本输出
	Err         string             `json:"err" bson:"err"`               // 错误原因
	PlanTime    int64              `json:"plan_time" bson:"plan_time"`   // 理论上的调度时间
	RealTime    int64              `json:"real_time" bson:"real_time"`   // 实际的调度时间
	StartTime   int64              `json:"start_time" bson:"start_time"` // 启动时间
	EndTime     int64              `json:"end_time" bson:"end_time"`     // 结束时间
}
