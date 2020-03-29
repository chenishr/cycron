package api

import (
	"cycron/libs"
	"cycron/mod"
	"cycron/sched"
	"encoding/json"
	"fmt"
	"github.com/gorhill/cronexpr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

func doDelTask(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		taskId  int
		postId  string
		bytes   []byte
		delCond bson.M
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "content-type")

	fmt.Println("coming a ", req.Method, " request: ", time.Now())
	if "OPTIONS" == req.Method {
		err = ServerError("忽略 OPTIONS 请求")
		goto ERR
	}

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postId = req.PostForm.Get("taskId")
	if "" == postId {
		err = ServerError("请求参数错误")
		goto ERR
	}

	taskId, _ = strconv.Atoi(postId)

	// 删除对应的 task
	delCond = bson.M{"_id": taskId}
	err = mod.GTaskMgr.DelTasks(delCond)

	sched.GScheduler.RemoveJob(int64(taskId))

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func doUpdateStatus(resp http.ResponseWriter, req *http.Request) {
	var (
		err        error
		taskId     int
		taskStatus int
		postId     string
		postStatus string
		bytes      []byte
		task       *mod.TaskMod
		findCond   bson.M
		uptCond    bson.M
		uptData    bson.M
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "content-type")

	fmt.Println("coming a ", req.Method, " request: ", time.Now())
	if "OPTIONS" == req.Method {
		err = ServerError("忽略 OPTIONS 请求")
		goto ERR
	}

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postId = req.PostForm.Get("taskId")
	postStatus = req.PostForm.Get("taskStatus")
	if "" == postId || "" == postStatus {
		err = ServerError("请求参数错误")
		goto ERR
	}

	taskId, _ = strconv.Atoi(postId)
	taskStatus, _ = strconv.Atoi(postStatus)

	// 获取对应的 task
	findCond = bson.M{"_id": taskId}
	if task, err = mod.GTaskMgr.FindOneTask(findCond); err != nil {
		goto ERR
	}

	if task.Status == taskStatus {
		err = ServerError("任务状态更新失败")
		goto ERR
	}

	// 更新任务状态
	uptCond = bson.M{"_id": task.Id}
	uptData = bson.M{
		"$set": bson.M{
			"status": taskStatus,
		},
	}
	if err = mod.GTaskMgr.UpdateOne(uptCond, uptData); err != nil {
		goto ERR
	}

	if 0 == taskStatus {
		// 移除任务
		sched.GScheduler.RemoveJob(task.Id)
	} else {
		// 添加任务
		sched.GScheduler.AddJob(task, false)
	}

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func doRunTask(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		taskId   int
		postId   string
		bytes    []byte
		task     *mod.TaskMod
		findCond bson.M
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "content-type")

	fmt.Println("coming a ", req.Method, " request: ", time.Now())
	if "OPTIONS" == req.Method {
		err = ServerError("忽略 OPTIONS 请求")
		goto ERR
	}

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postId = req.PostForm.Get("taskId")
	if "" == postId {
		err = ServerError("任务 ID 不能为空")
		goto ERR
	}

	taskId, _ = strconv.Atoi(postId)

	// 获取对应的 task
	findCond = bson.M{"_id": taskId}
	if task, err = mod.GTaskMgr.FindOneTask(findCond); err != nil {
		goto ERR
	}

	fmt.Println(task)

	// 执行任务
	sched.GScheduler.RunOnce(task.Id)

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func doSaveTask(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		postTask string
		bytes    []byte
		task     mod.TaskMod
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "content-type")

	fmt.Println("coming a ", req.Method, " request: ", time.Now())
	if "OPTIONS" == req.Method {
		err = ServerError("忽略 OPTIONS 请求")
		goto ERR
	}

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postTask = req.PostForm.Get("task")
	fmt.Println(string(bytes))

	// 3, 反序列化task
	if err = json.Unmarshal([]byte(postTask), &task); err != nil {
		goto ERR
	}
	// 4, 保存到mongoDB
	if 0 == task.Timeout {
		task.Timeout = 86400
	}
	if task.Concurrent < 1 {
		task.Concurrent = 1
	}

	if _, err = cronexpr.Parse(task.CronSpec); err != nil {
		err = ServerError("cron 表达式错误：" + err.Error())
		goto ERR
	}

	if "" == task.Command {
		err = ServerError("命令错误")
		goto ERR
	}

	if err = mod.GTaskMgr.UpsertDoc(&task); err != nil {
		goto ERR
	}

	fmt.Println(task)

	if err = sched.GScheduler.AddJob(&task, true); err != nil {
		fmt.Println("任务加入调度队列失败", err)
	}

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 列举所有crontab任务
func listTask(resp http.ResponseWriter, req *http.Request) {
	var (
		tasks []*mod.TaskMod
		bytes []byte
		err   error
		list  []map[string]interface{}
		expr  *cronexpr.Expression
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// 获取任务列表
	findCond := primitive.M{}
	if tasks, err = mod.GTaskMgr.FindTasks(findCond); err != nil {
		goto ERR
	}

	list = make([]map[string]interface{}, len(tasks))
	for k, v := range tasks {
		row := make(map[string]interface{})
		row["Id"] = v.Id
		row["TaskName"] = v.TaskName
		row["CronSpec"] = v.CronSpec
		row["Command"] = v.Command
		row["Status"] = v.Status
		row["GroupId"] = v.GroupId
		row["Concurrent"] = v.Concurrent
		row["Description"] = v.Description
		row["Notify"] = v.Notify
		row["NotifyEmail"] = v.NotifyEmail
		row["Timeout"] = v.Timeout

		if v.PrevTime == "" {
			row["PrevTime"] = "-"
		} else {
			row["PrevTime"] = v.PrevTime
		}

		if v.Status == 1 {
			if expr, err = cronexpr.Parse(v.CronSpec); err != nil {
				row["NextTime"] = "-"
			} else {
				row["NextTime"] = expr.Next(time.Now()).Format("2006-01-02 15:04:05")
			}
		} else {
			row["NextTime"] = "-"
		}

		list[k] = row
	}

	// 正常应答
	if bytes, err = libs.BuildResponse(0, "success", list); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}
