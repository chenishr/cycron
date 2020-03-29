package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"github.com/simagix/keyhole/mdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
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
	PlanTime    string `bson:"plan_time"`  // 理论上的调度时间
	RealTime    string `bson:"real_time"`  // 实际的调度时间
	StartTime   string `bson:"start_time"` // 启动时间
	EndTime     string `bson:"end_time"`   // 结束时间
}

type StatGroup struct {
	Day    string `bson:"day"`
	Status int    `bson:"status"`
}

type StatRes struct {
	Group StatGroup `bson:"_id"`
	Count float64   `bson:"count"`
}

type TaskLogMgr struct {
}

var (
	GTaskLogMgr *TaskLogMgr
)

func init() {
	GTaskLogMgr = &TaskLogMgr{}
}

func (tlm *TaskLogMgr) LogStat() (res []StatRes, err error) {
	var cur *mongo.Cursor

	today := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	pipeline := `
		[
		  {
			"$match": {
			  "start_time": { "$gt": "` + today + `" }
			}
		  },
		  {
			"$project": {
			  "status": 1,
			  "day": { "$substr": [ "$start_time", 0, 10 ] }
			}
		  },
		  {
			"$group": {
			  "_id": { "day": "$day", "status": "$status" },
			  "count": { "$sum": 1 }
			}
		  },
			{"$sort":{"_id.day":1}}
		]
		`
	fmt.Println(pipeline)

	opts := options.Aggregate()
	collection := dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLog)
	if cur, err = collection.Aggregate(context.TODO(), mdb.MongoPipeline(pipeline), opts); err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &res); err != nil {
		return nil, err
	}

	return
}

func (tm *TaskLogMgr) FindOneTaskLog(findCond interface{}) (taskLog *TaskLogMod, err error) {
	var (
		collection  *mongo.Collection
		res         *mongo.SingleResult
		findOptions *options.FindOneOptions
		findTaskLog TaskLogMod
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLog)

	findOptions = options.FindOne()
	res = collection.FindOne(context.TODO(), findCond, findOptions)

	if err = res.Decode(&findTaskLog); err != nil {
		return nil, err
	}

	taskLog = &findTaskLog
	return
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
