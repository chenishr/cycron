package handle

import (
	error2 "cycron/api/error"
	"cycron/libs"
	"cycron/mod"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func LogStat(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		bytes   []byte
		logData map[string]interface{}
		logStat []mod.StatRes
		list    []map[string]interface{}
	)

	// 删除对应的 task
	if logStat, err = mod.GTaskLogStatMgr.LogStat(); err != nil {
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
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func LogDetail(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		logId    int
		postId   string
		bytes    []byte
		findCond bson.M
		taskLog  *mod.TaskLogMod
		logData  map[string]interface{}
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	postId = req.PostForm.Get("logId")
	if "" == postId {
		err = error2.ServerError("请求参数错误")
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

	if taskLog.Status != mod.TASK_IGNORE {
		logData["StartTime"] = taskLog.StartTime
		logData["EndTime"] = taskLog.EndTime
	} else {
		logData["StartTime"] = "-"
		logData["EndTime"] = "-"
	}

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", logData); err == nil {
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
func ListLogs(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		taskId   int
		postId   string
		page     int
		pageSize int
		count    int64
		bytes    []byte
		findCond bson.M
		taskLogs []*mod.TaskLogMod
		list     []map[string]interface{}
		resJson  map[string]interface{}
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	postId = req.PostForm.Get("taskId")
	if "" == postId {
		err = error2.ServerError("请求参数错误")
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
	if taskLogs, err = mod.GTaskLogMgr.FindTaskLogs(findCond, int64(page), int64(pageSize)); err != nil {
		goto ERR
	}
	if count, err = mod.GTaskLogMgr.Count(findCond); err != nil {
		goto ERR
	}

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

	resJson = make(map[string]interface{})
	resJson["list"] = list
	resJson["page"] = page
	resJson["page_size"] = pageSize
	resJson["count"] = count
	resJson["total_page"] = math.Ceil(float64(count) / float64(pageSize))

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", resJson); err == nil {
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
