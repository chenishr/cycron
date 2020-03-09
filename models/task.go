package models

type TaskMod struct {
	Id           	int
	UserId       	int
	GroupId      	int
	TaskName     	string
	TaskType     	int
	Description  	string
	CronSpec     	string
	Concurrent   	int
	Command      	string
	Status       	int
	Notify       	int
	NotifyEmail  	string
	Timeout      	int
	ExecuteTimes 	int
	PrevTime     	int64
	CreateTime   	int64
}

/**
	获取待执行的任务
 */
func (t *TaskMod)GetTasks() (Tasks []*TaskMod,err error) {

	return
}