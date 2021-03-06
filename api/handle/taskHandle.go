package handle

import (
	error2 "cycron/api/error"
	"cycron/libs"
	"cycron/mod"
	"cycron/sched"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"net/http"
	"strconv"
	"time"
)

func DoDelTask(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		taskId  int
		postId  string
		bytes   []byte
		delCond bson.M
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postId = req.PostForm.Get("taskId")
	if "" == postId {
		err = error2.ServerError("请求参数错误")
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
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func DoUpdateStatus(resp http.ResponseWriter, req *http.Request) {
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

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postId = req.PostForm.Get("taskId")
	postStatus = req.PostForm.Get("taskStatus")
	if "" == postId || "" == postStatus {
		err = error2.ServerError("请求参数错误")
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
		err = error2.ServerError("任务状态更新失败")
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
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func DoRunTask(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		taskId   int
		postId   string
		bytes    []byte
		task     *mod.TaskMod
		findCond bson.M
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postId = req.PostForm.Get("taskId")
	if "" == postId {
		err = error2.ServerError("任务 ID 不能为空")
		goto ERR
	}

	taskId, _ = strconv.Atoi(postId)

	// 获取对应的 task
	findCond = bson.M{"_id": taskId}
	if task, err = mod.GTaskMgr.FindOneTask(findCond); err != nil {
		goto ERR
	}

	// 执行任务
	sched.GScheduler.RunOnce(task.Id)

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func DoSaveTask(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		postTask string
		bytes    []byte
		task     mod.TaskMod
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postTask = req.PostForm.Get("task")

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
		err = error2.ServerError("cron 表达式错误：" + err.Error())
		goto ERR
	}

	if "" == task.Command {
		err = error2.ServerError("命令错误")
		goto ERR
	}

	if err = mod.GTaskMgr.UpsertDoc(&task); err != nil {
		goto ERR
	}

	if err = sched.GScheduler.AddJob(&task, true); err != nil {
		log.Errorln("任务加入调度队列失败", err)
	}

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 列举所有crontab任务
func ListTask(resp http.ResponseWriter, req *http.Request) {
	var (
		postId   string
		page     int
		pageSize int
		tasks    []*mod.TaskMod
		bytes    []byte
		err      error
		list     []map[string]interface{}
		expr     *cronexpr.Expression
		findCond bson.M
		count    int64
		resJson  map[string]interface{}
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	postId = req.PostForm.Get("page")
	if "" == postId {
		page = 1
	} else {
		page, _ = strconv.Atoi(postId)
	}

	postId = req.PostForm.Get("page_size")
	if "" == postId {
		pageSize = 10
	} else {
		pageSize, _ = strconv.Atoi(postId)
	}

	// 获取任务列表
	findCond = primitive.M{}
	if tasks, err = mod.GTaskMgr.FindTasks(findCond, int64(page), int64(pageSize)); err != nil {
		goto ERR
	}
	if count, err = mod.GTaskMgr.Count(findCond); err != nil {
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

	resJson = make(map[string]interface{})
	resJson["list"] = list
	resJson["page"] = page
	resJson["page_size"] = pageSize
	resJson["count"] = count
	resJson["total_page"] = math.Ceil(float64(count) / float64(pageSize))

	// 正常应答
	if bytes, err = libs.BuildResponse(0, "success", resJson); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}
