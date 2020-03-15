package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskGroupMod struct {
	Id          primitive.ObjectID
	UserId      primitive.ObjectID
	GroupName   string
	Description string
	CreateTime  int64
	UpdateTime  int64
}

type TaskGroupMgr struct {
}

var (
	GTaskGroupMgr *TaskGroupMgr
)

func init() {
	GTaskGroupMgr = &TaskGroupMgr{}
}

func (g *TaskGroupMgr) AddGroup(taskGroup *TaskGroupMod) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertOneResult
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskGroup)

	if result, err = collection.InsertOne(context.TODO(), taskGroup); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 插入的 ID：", result)
	return
}
