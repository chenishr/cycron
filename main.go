package main

import (
	"cycron/api"
	"cycron/conf"
	"cycron/libs"
	"cycron/mod"
	"cycron/sched"
	"fmt"
	"time"
)

func main() {
	var (
		err error
	)
	fmt.Println("这是一个定时任务管理程序")

	//addTask()
	//addUser()
	//addGroup()

	// 启动 HTTPServer
	err = api.InitHttpServer()
	if err != nil {
		fmt.Println("HttpServer启动失败：", err)
	}

	err = sched.GScheduler.InitScheduler()
	if err != nil {
		fmt.Println("调度器启动失败：", err)
	}

	sched.GScheduler.Print()

	time.Sleep(3600 * time.Second)
	fmt.Println("一小时后程序退出")
}

func addGroup() {
	id, err := mod.GCommonMgr.GetMaxId(conf.GConfig.Models.TaskGroup)
	if err != nil {
		fmt.Println("获取自增 ID 失败", err)
	}

	taskGroup := &mod.TaskGroupMod{
		Id:          id,
		UserId:      0,
		GroupName:   "默认分组",
		Description: "默认分组",
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	err = mod.GTaskGroupMgr.AddGroup(taskGroup)
	if err != nil {
		fmt.Println("添加任务分组失败：", err)
	}
}

func addUser() {
	id, err := mod.GCommonMgr.GetMaxId(conf.GConfig.Models.User)
	if err != nil {
		fmt.Println("获取自增 ID 失败", err)
	}

	user := &mod.UserMod{
		Id:            id,
		UserName:      "chenishr",
		Password:      "",
		Email:         "chenishr@gmail.com",
		LastLoginTime: 0,
		LastIp:        "",
		Status:        1,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	user.Password = libs.Md5("123456" + string(user.CreateTime))

	err = mod.GUserMgr.AddUser(user)
	if err != nil {
		fmt.Println("添加用户失败：", err)
	}
}

func addTask() {
	var (
		err error
		id  int64
	)

	id, err = mod.GCommonMgr.GetMaxId(conf.GConfig.Models.Task)
	if err != nil {
		fmt.Println("获取自增 ID 失败", err)
	}

	task := &mod.TaskMod{
		Id:           id,
		UserId:       0,
		GroupId:      0,
		TaskName:     "Hello,Cycron!",
		TaskType:     0,
		Description:  "只是显示一个Cycron!",
		CronSpec:     "1/3 * * * * * *",
		Concurrent:   1,
		Command:      "echo 'Hello,Cycron!';",
		Status:       1,
		Notify:       1,
		NotifyEmail:  "chenishr@163.com",
		Timeout:      0,
		ExecuteTimes: 0,
		PrevTime:     0,
		CreateTime:   time.Now().Unix(),
	}

	err = mod.GTaskMgr.AddTask(task)
	if err != nil {
		fmt.Println("添加任务失败：", err)
	}
}
