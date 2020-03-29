package api

import (
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
	mux.HandleFunc("/task/save", doSaveTask)
	mux.HandleFunc("/task/run", doRunTask)
	mux.HandleFunc("/task/del", doDelTask)
	mux.HandleFunc("/task/update_status", doUpdateStatus)
	mux.HandleFunc("/task/list", listTask)

	// 日记管理
	mux.HandleFunc("/log/list", listLogs)
	mux.HandleFunc("/log/detail", logDetail)
	mux.HandleFunc("/log/stat", logStat)

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
