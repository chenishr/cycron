package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// 任务执行日志
type TaskLogMod struct {
	Id          int64  `bson:"_id"`
	TaskId      int64  `bson:"task_id"`
	Command     string `bson:"command"` // 脚本命令
	Status      int    `bson:"status"`
	ProcessTime int    `bson:"process_time"`
	Output      string `bson:"output"`     // 脚本输出
	Err         string `bson:"err"`        // 错误原因
	PlanTime    int64  `bson:"plan_time"`  // 理论上的调度时间
	RealTime    int64  `bson:"real_time"`  // 实际的调度时间
	StartTime   int64  `bson:"start_time"` // 启动时间
	EndTime     int64  `bson:"end_time"`   // 结束时间
}

type TaskLogMgr struct {
}

var (
	GTaskLogMgr *TaskLogMgr
)

func init() {
	GTaskLogMgr = &TaskLogMgr{}
}

func (tlm *TaskLogMgr) FindTaskLogs(findCond interface{}, page, pageSize int64) (taskLogs []*TaskLogMod, err error) {
	var (
		collection  *mongo.Collection
		cursor      *mongo.Cursor
		findOptions *options.FindOptions
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLog)

	if page < 1 {
		page = 1
	}

	findOptions = options.Find()
	findOptions.SetSort(bsonx.Doc{{"_id", bsonx.Int32(-1)}})
	findOptions.SetSkip((page - 1) * pageSize)
	findOptions.SetLimit(pageSize)
	cursor, err = collection.Find(context.TODO(), findCond, findOptions)
	if err != nil {
		return nil, err
	}

	// 遍历获取所有的文档
	for cursor.Next(context.TODO()) {
		var taskLog TaskLogMod
		// decode the document into given type
		if err := cursor.Decode(&taskLog); err != nil {
			return nil, err
		}
		taskLogs = append(taskLogs, &taskLog)
	}
	return
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

	fmt.Println(" 插入的 task_log ID：", result)
	return
}
