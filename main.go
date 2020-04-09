package main

import (
	"cycron/api"
	"cycron/conf"
	"cycron/dbs"
	"cycron/mod"
	"cycron/sched"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	var (
		file *os.File
		err  error
	)
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if file, err = os.OpenFile(conf.GConfig.Logger.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm); err != nil {
		log.Fatalln("打开日志文件错误：", err)
	} else {
		log.SetOutput(file)
	}

	// 设置日志级别为warn以上
	log.SetLevel(log.DebugLevel)
	log.SetLevel(log.Level(conf.GConfig.Logger.LogLevel))

	// 初始化 MongoDB 连接池
	dbs.InitMongoPool()

	// 初始化任务调度器
	sched.InitScheduler()

	// 初始化默认用户
	mod.GUserMgr.InitUser()

	//mod.GTaskLogStatMgr.InitData()
}

func main() {
	var (
		err error
	)
	log.Info("这是一个定时任务管理程序")

	// 启动 HTTPServer
	err = api.InitHttpServer()
	if err != nil {
		log.Fatalln("HttpServer启动失败：", err)
	}

	// 启动任务调度器
	err = sched.GScheduler.StartScheduler()
	if err != nil {
		log.Fatalln("调度器启动失败：", err)
	}
}
