package handle

import (
	error2 "cycron/api/error"
	"cycron/libs"
	"cycron/mod"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"net/http"
	"strconv"
)

func DoSaveTaskGroup(resp http.ResponseWriter, req *http.Request) {
	var (
		err       error
		postTask  string
		bytes     []byte
		taskGroup mod.TaskGroupMod
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	//postTask, _ = ioutil.ReadAll(req.Body)
	postTask = req.PostForm.Get("taskGroup")

	// 3, 反序列化task
	if err = json.Unmarshal([]byte(postTask), &taskGroup); err != nil {
		goto ERR
	}
	// 4, 保存到mongoDB

	if "" == taskGroup.GroupName {
		err = error2.ServerError("参数错误")
		goto ERR
	}

	if err = mod.GTaskGroupMgr.UpsertDoc(&taskGroup); err != nil {
		goto ERR
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
func ListTaskGroup(resp http.ResponseWriter, req *http.Request) {
	var (
		postId     string
		page       int
		pageSize   int
		taskGroups []*mod.TaskGroupMod
		bytes      []byte
		err        error
		list       []map[string]interface{}
		count      int64
		resJson    map[string]interface{}
		findCond   primitive.M
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
	if taskGroups, err = mod.GTaskGroupMgr.FindTaskGroups(findCond, int64(page), int64(pageSize)); err != nil {
		goto ERR
	}
	if count, err = mod.GTaskGroupMgr.Count(findCond); err != nil {
		goto ERR
	}

	list = make([]map[string]interface{}, len(taskGroups))
	for k, v := range taskGroups {
		row := make(map[string]interface{})
		row["Id"] = v.Id
		row["UserId"] = v.UserId
		row["GroupName"] = v.GroupName
		row["Description"] = v.Description

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
