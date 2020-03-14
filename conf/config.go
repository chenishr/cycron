package conf

import (
	"encoding/json"
	"io/ioutil"
)

// 程序配置
type Config struct {
	Mail MailConf 	`json:"mail"`
	Mongo MongoConf `json:"mongo"`
	Logger LoggerConf `json:"logger"`
	Models ModelsConf	`json:"models"`
}

type MailConf struct {
	User 		string 	`json:"user"`
	Host 		string 	`json:"host"`
	PassWord	string 	`json:"password"`
	Port		int 	`json:"port"`
}

type MongoConf struct {
	Uri 			string 	`json:"uri"`
	ConnectTimeout 	int 	`json:"connectTimeout"`
	CommitTimeout 	int 	`json:"commitTimeout"`
}

type LoggerConf struct {
	BatchSize 		int 	`json:"batchsize"`
}

type ModelsConf struct {
	Db 				string 	`json:"db"`
	TaskLog 		string 	`json:"taskLog"`
	Task 			string 	`json:"task"`
}

var (
	// 单例
	GConfig *Config
)

// 加载配置
func init() {
	var (
		content 	[]byte
		config 		Config
		filename	string
		err			error
	)

	filename = "./conf/local.config.json"

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
