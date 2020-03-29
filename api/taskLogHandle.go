package api

import (
	"cycron/libs"
	"cycron/mod"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func logStat(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		bytes   []byte
		logData map[string]interface{}
		logStat []mod.StatRes
		list    []map[string]interface{}
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

	// 删除对应的 task
	if logStat, err = mod.GTaskLogMgr.LogStat(); err != nil {
		goto ERR
	}

	list = make([]map[string]interface{}, len(logStat))
	for i, v := range logStat {
		logData = make(map[string]interface{})
		logData["day"] = v.Group.Day
		logData["status"] = v.Group.Status
		logData["count"] = v.Count

		list[i] = logData
	}

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", list); err == nil {
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

func logDetail(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		logId    int
		postId   string
		bytes    []byte
		findCond bson.M
		taskLog  *mod.TaskLogMod
		logData  map[string]interface{}
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

	postId = req.PostForm.Get("logId")
	if "" == postId {
		err = ServerError("请求参数错误")
		goto ERR
	}
	logId, _ = strconv.Atoi(postId)

	// 删除对应的 task
	findCond = bson.M{"_id": logId}
	if taskLog, err = mod.GTaskLogMgr.FindOneTaskLog(findCond); err != nil {
		goto ERR
	}

	logData = make(map[string]interface{})
	logData["Id"] = taskLog.Id
	logData["TaskId"] = taskLog.TaskId
	logData["Command"] = taskLog.Command
	logData["Status"] = taskLog.Status
	logData["Output"] = strings.Trim(taskLog.Output, " \n")
	logData["Err"] = strings.Trim(taskLog.Err, " \n")
	logData["ProcessTime"] = taskLog.ProcessTime
	logData["PlanTime"] = taskLog.PlanTime
	logData["RealTime"] = taskLog.RealTime
	logData["StartTime"] = taskLog.StartTime
	logData["EndTime"] = taskLog.EndTime

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", logData); err == nil {
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
func listLogs(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		taskId   int
		postId   string
		page     int
		pageSize int
		bytes    []byte
		findCond bson.M
		taskLogs []*mod.TaskLogMod
		list     []map[string]interface{}
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

	postId = req.PostForm.Get("taskId")
	if "" == postId {
		err = ServerError("请求参数错误")
		goto ERR
	}
	taskId, _ = strconv.Atoi(postId)

	postId = req.PostForm.Get("page")
	if "" == postId {
		page = 1
	} else {
		page, _ = strconv.Atoi(postId)
	}

	postId = req.PostForm.Get("page_size")
	if "" == postId {
		pageSize = 20
	} else {
		pageSize, _ = strconv.Atoi(postId)
	}

	// 删除对应的 task
	findCond = bson.M{"task_id": taskId}
	taskLogs, err = mod.GTaskLogMgr.FindTaskLogs(findCond, int64(page), int64(pageSize))

	list = make([]map[string]interface{}, len(taskLogs))
	for k, v := range taskLogs {
		row := make(map[string]interface{})
		row["Id"] = v.Id
		row["TaskId"] = v.TaskId
		row["Status"] = v.Status
		row["ProcessTime"] = v.ProcessTime
		row["RealTime"] = v.RealTime

		list[k] = row
	}

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", list); err == nil {
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
