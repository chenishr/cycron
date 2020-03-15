package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskLogMgr struct {
}

var (
	GTaskLogMgr *TaskLogMgr
)

func init() {
	GTaskLogMgr = &TaskLogMgr{}
}

func (tlm *TaskLogMgr) InsertMany(taskLog []interface{}) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertManyResult
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLog)

	if result, err = collection.InsertMany(context.TODO(), taskLog); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 插入的 ID：", result)
	return
}
