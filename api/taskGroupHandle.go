package api

import (
	"cycron/libs"
	"cycron/mod"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func doSaveTaskGroup(resp http.ResponseWriter, req *http.Request) {
	var (
		err       error
		postTask  string
		bytes     []byte
		taskGroup mod.TaskGroupMod
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
	postTask = req.PostForm.Get("taskGroup")
	fmt.Println(string(bytes))

	// 3, 反序列化task
	if err = json.Unmarshal([]byte(postTask), &taskGroup); err != nil {
		goto ERR
	}
	// 4, 保存到mongoDB

	if "" == taskGroup.GroupName {
		err = ServerError("参数错误")
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
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 列举所有crontab任务
func listTaskGroup(resp http.ResponseWriter, req *http.Request) {
	var (
		taskGroups []*mod.TaskGroupMod
		bytes      []byte
		err        error
		list       []map[string]interface{}
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// 获取任务列表
	findCond := primitive.M{}
	if taskGroups, err = mod.GTaskGroupMgr.FindTaskGroups(findCond); err != nil {
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
