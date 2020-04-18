package dbs

import (
	"context"
	"cycron/conf"
	"github.com/silenceper/pool"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	GMongoPool pool.Pool
)

func InitMongoPool() {
	//创建一个连接池： 初始化5，最大链接30
	//创建一个连接池： 初始化5，最大空闲连接是20，最大并发连接30
	poolConfig := &pool.Config{
		InitialCap: 5,  //资源池初始连接数
		MaxIdle:    20, //最大空闲连接数
		MaxCap:     30, //最大并发连接数
		Factory:    Factory,
		Close:      Close,
		Ping:       Ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		log.Info("err=", err)

	}

	GMongoPool = p
}

func Factory() (interface{}, error) {
	var (
		client    *mongo.Client
		mongoConf conf.MongoConf
		ctx       context.Context
		err       error
	)

	mongoConf = conf.GConfig.Mongo

	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(mongoConf.ConnectTimeout)*time.Millisecond)

	// 建立mongodb连接
	log.Info("建立mongodb连接")
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoConf.Uri)); err != nil {
		log.Fatalln("链接 MongoDB 失败：", err)
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalln("ping MongoDB 失败：", err)
		return nil, err
	}

	return client, nil
}

func Close(v interface{}) error {
	return v.(*mongo.Client).Disconnect(context.TODO())
}

func Ping(v interface{}) error {
	return v.(*mongo.Client).Ping(context.TODO(), nil)
}
