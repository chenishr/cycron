package api

import (
	"cycron/conf"
	"cycron/libs"
	"cycron/mod"
	"cycron/sched"
	"encoding/json"
	"fmt"
	"github.com/gorhill/cronexpr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 任务的HTTP接口
type HttpServer struct {
	Server *http.Server
}

var (
	// 单例对象
	GServer *HttpServer
)

type ServerError string

func (e ServerError) Error() string {
	return "httpserver error: " + strconv.Quote(string(e))
}

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
		row["Status"] = v.Status
		row["Description"] = v.Description
		row["Notify"] = v.Notify
		row["NotifyEmail"] = v.NotifyEmail
		row["Timeout"] = v.Timeout

		if v.PrevTime > 0 {
			row["PrevTime"] = time.Unix(v.PrevTime, 0).Format("2006-01-02 15:04:05")
		} else {
			row["PrevTime"] = "-"
		}

		if v.Status == 1 {
			if expr, err = cronexpr.Parse(v.CronSpec); err != nil {
				row["NextTime"] = ""
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
		row["RealTime"] = time.Unix(v.RealTime, 0).Format("2006-01-02 15:04:05")

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

// 初始化服务
func InitHttpServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		Server        *http.Server
		staticDir     http.Dir     // 静态文件根目录
		staticHandler http.Handler // 静态文件的HTTP回调
		serverConf    conf.ServerConf
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/task/save", doSaveTask)
	mux.HandleFunc("/task/run", doRunTask)
	mux.HandleFunc("/task/del", doDelTask)
	mux.HandleFunc("/task/update_status", doUpdateStatus)
	mux.HandleFunc("/task/list", listTask)
	mux.HandleFunc("/log/list", listLogs)

	//  /index.html
	serverConf = conf.GConfig.Server

	// 静态文件目录
	staticDir = http.Dir(serverConf.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler)) //   ./webroot/index.html

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(serverConf.Port)); err != nil {
		return
	}

	// 创建一个HTTP服务
	Server = &http.Server{
		ReadTimeout:  time.Duration(serverConf.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(serverConf.WriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 赋值单例
	GServer = &HttpServer{
		Server: Server,
	}

	// 启动了服务端
	go Server.Serve(listener)

	return
}
