package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommonMod struct {
	Id    string `bson:"_id"`
	MaxId int64  `bson:"max_id"`
}

type CommonMgr struct {
}

var (
	GCommonMgr *CommonMgr
)

func init() {
	GCommonMgr = &CommonMgr{}
}

func (tm *CommonMgr) GetMaxId(coll string) (maxId int64, err error) {
	var (
		collection *mongo.Collection
		res        *mongo.SingleResult
		findCond   interface{}
		uptData    interface{}
		uptOptions *options.FindOneAndUpdateOptions
		doc        CommonMod
		client     *mongo.Client
		p          interface{}
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Common)

	findCond = bson.M{"_id": coll}
	uptData = bson.M{"$inc": bson.M{"max_id": 1}}
	uptOptions = &options.FindOneAndUpdateOptions{}
	uptOptions.SetUpsert(true)
	uptOptions.SetReturnDocument(1)
	res = collection.FindOneAndUpdate(context.TODO(), findCond, uptData, uptOptions)

	err = res.Decode(&doc)
	if err != nil {
		return
	}

	return doc.MaxId, nil
}
