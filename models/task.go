package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// 任务完成状态
const (
	TASK_SUCCESS 	= 0  	// 任务执行成功
	TASK_ERROR   	= 1 	// 任务执行出错
	TASK_TIMEOUT 	= 2 	// 任务执行超时
	TASK_CANCEL 	= 3 	// 任务被取消
)

type TaskMod struct {
	Id           	primitive.ObjectID 	`bson:"_id"`
	UserId       	int					`bson:"user_id"`
	GroupId      	int					`bson:"group_id"`
	TaskName     	string				`bson:"task_name"`
	TaskType     	int					`bson:"task_type"`
	Description  	string				`bson:"description"`
	CronSpec     	string				`bson:"cron_spec"`
	Concurrent   	int					`bson:"concurrent"`
	Command      	string				`bson:"command"`
	Status       	int					`bson:"status"`					// 0 停止；1 启动
	Notify       	int					`bson:"notify"`					// 0 不通知；1 执行失败通知；2 执行结束通知
	NotifyEmail  	string				`bson:"notify_email"`
	Timeout      	int					`bson:"timeout"`
	ExecuteTimes 	int					`bson:"execute_times"`
	PrevTime     	int64				`bson:"prev_time"`
	CreateTime   	int64				`bson:"create_time"`
}
