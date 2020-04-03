package api

import (
	"cycron/api/handle"
	mw "cycron/api/middleware"
	"cycron/conf"
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
	// 任务管理
	mux.HandleFunc("/task/save", mw.Chain(handle.DoSaveTask, mw.CheckToken(), mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/task/run", mw.Chain(handle.DoRunTask, mw.CheckToken(), mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/task/del", mw.Chain(handle.DoDelTask, mw.CheckToken(), mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/task/update_status", mw.Chain(handle.DoUpdateStatus, mw.CheckToken(), mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/task/list", mw.Chain(handle.ListTask, mw.ReturnOption(), mw.SetHeader()))

	// 日记管理
	mux.HandleFunc("/log/list", mw.Chain(handle.ListLogs, mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/log/detail", mw.Chain(handle.LogDetail, mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/log/stat", mw.Chain(handle.LogStat, mw.ReturnOption(), mw.SetHeader()))

	// 任务组管理
	mux.HandleFunc("/group/save", mw.Chain(handle.DoSaveTaskGroup, mw.CheckToken(), mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/group/list", mw.Chain(handle.ListTaskGroup, mw.ReturnOption(), mw.SetHeader()))

	// 用户管理
	mux.HandleFunc("/user/login", mw.Chain(handle.DoLogin, mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/user/save", mw.Chain(handle.DoSaveUser, mw.ReturnOption(), mw.SetHeader()))
	mux.HandleFunc("/user/info", mw.Chain(handle.UserInfo, mw.CheckToken(), mw.ReturnOption(), mw.SetHeader()))

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
