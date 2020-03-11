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
func GetTasks() (Tasks []*TaskMod,err error) {
	var (
		task *TaskMod
	)

	task = &TaskMod{
		Id:           1,
		UserId:       1,
		GroupId:      1,
		TaskName:     "第 1 个任务",
		TaskType:     0,
		Description:  "第 1 个任务",
		CronSpec:     "*/5 * * * * * *",
		Concurrent:   1,
		Command:      "echo 'Hello,World!';",
		Status:       1,
		Notify:       1,
		NotifyEmail:  "chenishr@163.com",
		Timeout:      0,
		ExecuteTimes: 0,
		PrevTime:     0,
		CreateTime:   0,
	}
	Tasks = append(Tasks,task)

	task = &TaskMod{
		Id:           2,
		UserId:       1,
		GroupId:      1,
		TaskName:     "第 2 个任务",
		TaskType:     0,
		Description:  "第 2 个任务",
		CronSpec:     "10 */5 * * * * *",
		Concurrent:   1,
		Command:      "echo 'Hello,golang!';",
		Status:       1,
		Notify:       0,
		NotifyEmail:  "",
		Timeout:      0,
		ExecuteTimes: 0,
		PrevTime:     0,
		CreateTime:   0,
	}
	Tasks = append(Tasks,task)

	task = &TaskMod{
		Id:           3,
		UserId:       1,
		GroupId:      1,
		TaskName:     "第 3 个任务",
		TaskType:     0,
		Description:  "第 3 个任务",
		CronSpec:     "*/10 * * * * * *",
		Concurrent:   1,
		Command:      "uptime",
		Status:       1,
		Notify:       1,
		NotifyEmail:  "chenishr@163.com\nchenishr@gmail.com",
		Timeout:      0,
		ExecuteTimes: 0,
		PrevTime:     0,
		CreateTime:   0,
	}
	Tasks = append(Tasks,task)

	return
}