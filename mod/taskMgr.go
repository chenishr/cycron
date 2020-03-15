package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskMgr struct {
}

var (
	GTaskMgr *TaskMgr
)

func init() {
	GTaskMgr = &TaskMgr{}
}

func (tm *TaskMgr) UpdateOne(uptCond interface{}, update interface{}) (err error) {
	var (
		res        *mongo.UpdateResult
		collection *mongo.Collection
	)

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

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	if result, err = collection.InsertOne(context.TODO(), task); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 插入的 ID：", result.InsertedID)
	return
}

func (tm *TaskMgr) FindTasks(findCond interface{}) (tasks []*TaskMod, err error) {
	var (
		collection *mongo.Collection
		cursor     *mongo.Cursor
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	cursor, err = collection.Find(context.TODO(), findCond)
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
