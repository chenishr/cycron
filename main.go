package main

import (
	"cycron/libs"
	"cycron/mod"
	"cycron/sched"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func main() {
	var(
		err error
	)
	fmt.Println("这是一个定时任务管理程序")

	//addTask()
	addUser()

	err = sched.GScheduler.InitScheduler()
	if err != nil {
		fmt.Println("调度器启动失败：",err)
	}

	sched.GScheduler.Print()
}

func addUser()  {
	user := &mod.UserMod{
		Id:            primitive.NewObjectID(),
		UserName:      "chenishr",
		Password:      "",
		Email:         "chenishr@gmail.com",
		LastLoginTime: 0,
		LastIp:        "",
		Status:        1,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	user.Password = libs.Md5("123456" + string(user.CreateTime));

	err := mod.GUserMgr.AddUser(user)
	if err != nil {
		fmt.Println("添加用户失败：",err)
	}
}

func addTask()  {
	var(
		err error
	)

	task := &mod.TaskMod{
		Id:           primitive.NewObjectID(),
		UserId:       primitive.ObjectID{},
		GroupId:      primitive.ObjectID{},
		TaskName:     "Hello,World!",
		TaskType:     0,
		Description:  "只是显示一个Cycron!",
		CronSpec:     "* * * * * * *",
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
		fmt.Println("添加任务失败：",err)
	}
}
