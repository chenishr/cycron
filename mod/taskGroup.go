package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type TaskGroupMod struct {
	Id          int64  `bson:"_id"`
	UserId      int64  `bson:"user_id"`
	GroupName   string `bson:"group_name"`
	Description string `bson:"description"`
	CreateTime  string `bson:"create_time"`
	UpdateTime  string `bson:"update_time"`
}

type TaskGroupMgr struct {
}

var (
	GTaskGroupMgr *TaskGroupMgr
)

func init() {
	GTaskGroupMgr = &TaskGroupMgr{}
}

/*
添加或者更新文档
*/
func (g *TaskGroupMgr) UpsertDoc(taskGroup *TaskGroupMod) (err error) {
	var (
		id      int64
		uptCond bson.M
		uptData bson.M
	)
	if taskGroup.Id == 0 {
		id, err = GCommonMgr.GetMaxId(conf.GConfig.Models.TaskGroup)

		taskGroup.Id = id
		taskGroup.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		taskGroup.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

		return g.AddGroup(taskGroup)
	}

	uptCond = bson.M{"_id": taskGroup.Id}
	uptData = bson.M{
		"$set": bson.M{
			"user_id":     taskGroup.UserId,
			"description": taskGroup.Description,
			"group_name":  taskGroup.GroupName,
			"update_time": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	return g.UpdateOne(uptCond, uptData)
}

func (g *TaskGroupMgr) UpdateOne(uptCond interface{}, update interface{}) (err error) {
	var (
		collection *mongo.Collection
		client     *mongo.Client
		p          interface{}
		result     *mongo.UpdateResult
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskGroup)

	// 执行删除
	if result, err = collection.UpdateOne(context.TODO(), uptCond, update); err != nil {
		log.Errorln(err)
		return
	}

	log.Debugln("更新"+conf.GConfig.Models.TaskGroup+"结果：", result)

	return
}

func (g *TaskGroupMgr) AddGroup(taskGroup *TaskGroupMod) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertOneResult
		client     *mongo.Client
		p          interface{}
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskGroup)

	if result, err = collection.InsertOne(context.TODO(), taskGroup); err != nil {
		log.Errorln(err)
		return
	}

	log.Debugln(" 插入的 task_group ID：", result)
	return
}

func (g *TaskGroupMgr) FindTaskGroups(findCond interface{}) (taskGroups []*TaskGroupMod, err error) {
	var (
		collection  *mongo.Collection
		cursor      *mongo.Cursor
		findOptions *options.FindOptions
		client      *mongo.Client
		p           interface{}
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskGroup)

	findOptions = options.Find()
	findOptions.SetSort(bsonx.Doc{{"_id", bsonx.Int32(-1)}})
	cursor, err = collection.Find(context.TODO(), findCond, findOptions)
	if err != nil {
		return nil, err
	}

	// 遍历获取所有的文档
	for cursor.Next(context.TODO()) {
		var taskGroup TaskGroupMod
		// decode the document into given type
		if err := cursor.Decode(&taskGroup); err != nil {
			return nil, err
		}
		taskGroups = append(taskGroups, &taskGroup)
	}
	return taskGroups, nil
}
