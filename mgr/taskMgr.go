package mgr

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"cycron/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskMgr struct {

}

var (
	GTaskMgr * TaskMgr
)

func init()  {
	GTaskMgr = &TaskMgr{}
}

/*
删除任务
 */
func (tm *TaskMgr)DelTasks(delCond interface{}) (err error) {
	var  (
		delResult *mongo.DeleteResult
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
func (tm *TaskMgr)AddTask(task *models.TaskMod) (err error) {
	var (
		collection *mongo.Collection
		result *mongo.InsertOneResult
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	if result, err = collection.InsertOne(context.TODO(), task); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 插入的 ID：",result.InsertedID)
	return
}

func (tm *TaskMgr)FindTasks(findCond interface{}) (tasks []*models.TaskMod, err error){
	var (
		collection *mongo.Collection
		cursor *mongo.Cursor
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.Task)

	cursor, err = collection.Find(context.TODO(), findCond)
	if err != nil {
		return nil, err
	}

	// iterate through all documents
	for cursor.Next(context.TODO()) {
		var task models.TaskMod
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
func (tm *TaskMgr)GetTasks() (tasks []*models.TaskMod,err error) {
	/*
	var (
		task *models.TaskMod
	)

	task = &models.TaskMod{
		Id:			  primitive.NewObjectID(),
		UserId:       1,
		GroupId:      1,
		TaskName:     "第 1 个任务",
		TaskType:     0,
		Description:  "第 1 个任务",
		CronSpec:     "5 * * * * * *",
		Concurrent:   1,
		Command:      "echo 'Hello,World!';",
		Status:       1,
		Notify:       1,
		NotifyEmail:  "chenishr@163.com",
		Timeout:      0,
		ExecuteTimes: 0,
		PrevTime:     0,
		CreateTime:   0,
	}

	err = tm.AddTask(task)
	 */

	findCond := primitive.M{"status":1}
	tasks,err = tm.FindTasks(findCond)

	return
}
