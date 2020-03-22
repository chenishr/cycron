package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskGroupMod struct {
	Id          int64  `bson:"_id"`
	UserId      int64  `bson:"user_id"`
	GroupName   string `bson:"group_name"`
	Description string `bson:"description"`
	CreateTime  int64  `bson:"create_time"`
	UpdateTime  int64  `bson:"update_time"`
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

	fmt.Println(" 插入的 task_group ID：", result)
	return
}
