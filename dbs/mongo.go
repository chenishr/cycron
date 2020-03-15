package dbs

import (
	"context"
	"cycron/conf"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	Client *mongo.Client
}

var (
	GMongo *Mongo
)

func init() {
	var (
		client    *mongo.Client
		mongoConf conf.MongoConf
		ctx       context.Context
		err       error
	)

	mongoConf = conf.GConfig.Mongo

	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(mongoConf.ConnectTimeout)*time.Millisecond)

	// 建立mongodb连接
	fmt.Println("建立mongodb连接")
	if client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoConf.Uri)); err != nil {
		fmt.Println("链接 MongoDB 失败：", err)
		return
	}
	fmt.Println("建立mongodb连接 成功")

	GMongo = &Mongo{Client: client}
	return
}
