package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// 任务完成状态
const (
	TASK_SUCCESS = 0 // 任务执行成功
	TASK_ERROR   = 1 // 任务执行出错
	TASK_TIMEOUT = 2 // 任务执行超时
	TASK_CANCEL  = 3 // 任务被取消
)

type TaskMod struct {
	Id           int64  `bson:"_id"`
	UserId       int64  `bson:"user_id"`
	GroupId      int64  `bson:"group_id"`
	TaskName     string `bson:"task_name"`
	TaskType     int    `bson:"task_type"`
	Description  string `bson:"description"`
	CronSpec     string `bson:"cron_spec"`
	Concurrent   int    `bson:"concurrent"`
	Command      string `bson:"command"`
	Status       int    `bson:"status"` // 0 停止；1 启动
	Notify       int    `bson:"notify"` // 0 不通知；1 执行失败通知；2 执行结束通知
	NotifyEmail  string `bson:"notify_email"`
	Timeout      int    `bson:"timeout"`
	ExecuteTimes int    `bson:"execute_times"`
	PrevTime     int64  `bson:"prev_time"`
	CreateTime   int64  `bson:"create_time"`
}

type TaskMgr struct {
}

var (
	GTaskMgr *TaskMgr
)

func init() {
	GTaskMgr = &TaskMgr{}
}

/*
添加或者更新文档
*/
func (tm *TaskMgr) UpsertDoc(task *TaskMod) (err error) {
	var (
		id      int64
		uptCond bson.M
		uptData bson.M
	)
	if task.Id == 0 {
		id, err = GCommonMgr.GetMaxId(conf.GConfig.Models.Task)
		task.Id = id

		//  默认启动任务
		task.Status = 1

		return tm.AddTask(task)
	}

	uptCond = bson.M{"_id": task.Id}
	uptData = bson.M{
		"$set": bson.M{
			"task_name":    task.TaskName,
			"description":  task.Description,
			"group_id":     task.GroupId,
			"concurrent":   task.Concurrent,
			"cron_spec":    task.CronSpec,
			"command":      task.Command,
			"timeout":      task.Timeout,
			"notify":       task.Notify,
			"notify_email": task.NotifyEmail,
		},
	}
	return tm.UpdateOne(uptCond, uptData)
}

func (tm *TaskMgr) UpdateOne(uptCond interface{}, update interface{}) (err error) {
	var (
		res        *mongo.UpdateResult
		collection *mongo.Collection
	)

	fmt.Println("执行更新")

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	// 执行删除
	if res, err = collection.UpdateOne(context.TODO(), uptCond, update); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)

	return
}

/*
删除任务
*/
func (tm *TaskMgr) DelTasks(delCond interface{}) (err error) {
	var (
		delResult  *mongo.DeleteResult
		collection *mongo.Collection
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	// 执行删除
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(delResult)

	return
}

/*
添加任务
*/
func (tm *TaskMgr) AddTask(task *TaskMod) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertOneResult
	)

	fmt.Println("执行添加")

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	if result, err = collection.InsertOne(context.TODO(), task); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 插入的 task ID：", result.InsertedID)
	return
}

func (tm *TaskMgr) FindOneTask(findCond interface{}) (task *TaskMod, err error) {
	var (
		collection  *mongo.Collection
		res         *mongo.SingleResult
		findOptions *options.FindOneOptions
		findTask    TaskMod
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	findOptions = options.FindOne()
	res = collection.FindOne(context.TODO(), findCond, findOptions)

	if err = res.Decode(&findTask); err != nil {
		return nil, err
	}

	task = &findTask
	return
}

func (tm *TaskMgr) FindTasks(findCond interface{}) (tasks []*TaskMod, err error) {
	var (
		collection  *mongo.Collection
		cursor      *mongo.Cursor
		findOptions *options.FindOptions
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	findOptions = options.Find()
	findOptions.SetSort(bsonx.Doc{{"_id", bsonx.Int32(-1)}})
	cursor, err = collection.Find(context.TODO(), findCond, findOptions)
	if err != nil {
		return nil, err
	}

	// 遍历获取所有的文档
	for cursor.Next(context.TODO()) {
		var task TaskMod
		// decode the document into given type
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

/**
获取待执行的任务
*/
func (tm *TaskMgr) GetTasks() (tasks []*TaskMod, err error) {

	findCond := primitive.M{"status": 1}
	tasks, err = tm.FindTasks(findCond)

	return
}
