package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 程序配置
type Config struct {
	Mail   MailConf   `json:"mail"`
	Mongo  MongoConf  `json:"mongo"`
	Logger LoggerConf `json:"logger"`
	Models ModelsConf `json:"models"`
	Server ServerConf `json:"server"`
}

type MailConf struct {
	User     string `json:"user"`
	Host     string `json:"host"`
	PassWord string `json:"password"`
	Port     int    `json:"port"`
}

type MongoConf struct {
	Uri            string `json:"uri"`
	ConnectTimeout int    `json:"connectTimeout"`
	CommitTimeout  int    `json:"commitTimeout"`
}

type LoggerConf struct {
	BatchSize int    `json:"batchsize"`
	LogFile   string `json:"logFile"`
	LogLevel  int    `json:"logLevel"`
}

type ModelsConf struct {
	Db          string `json:"db"`
	Task        string `json:"task"`
	TaskGroup   string `json:"taskGroup"`
	TaskLog     string `json:"taskLog"`
	TaskLogStat string `json:"taskLogStat"`
	User        string `json:"user"`
	Common      string `json:"common"`
}

type ServerConf struct {
	WebRoot      string `json:"webRoot"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
	JwtSecret    string `json:"jwtSecret"`
}

var (
	// 单例
	GConfig *Config
	RootPath string
)

// 加载配置
func init() {
	var (
		content  []byte
		config   Config
		filename string
		err      error
	)

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	RootPath = filepath.Dir(ex)

	filename = RootPath + "/conf/local.config.json"

	// 1, 把配置文件读进来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 2, 做JSON反序列化
	if err = json.Unmarshal(content, &config); err != nil {
		return
	}

	// 3, 赋值单例
	GConfig = &config

	return
}
