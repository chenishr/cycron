package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"github.com/simagix/keyhole/mdb"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 任务执行日志
type TaskLogStatMod struct {
	Id       int64  `bson:"_id"`
	Status   int    `bson:"status"`
	PlanTime string `bson:"plan_time"` // 理论上的调度时间
	Count    int64  `bson:"count"`     // 结束时间
}
type TaskLogStatMgr struct {
}

var (
	GTaskLogStatMgr *TaskLogStatMgr
)

func init() {
	GTaskLogStatMgr = &TaskLogStatMgr{}
}

/*
将日记信息同步到统计集合
*/
func (s *TaskLogStatMgr) InitData() {
	var (
		id       int64
		uptData  bson.M
		findCond bson.M
		err      error
		data     []StatRes
		statRes  StatRes
		stat     *TaskLogStatMod
	)

	data, err = GTaskLogMgr.LogStat(30, 13)
	if err != nil {
		log.Infoln("读取数据失败：", err)
		return
	}

	for _, statRes = range data {
		findCond = bson.M{
			"plan_time": statRes.Group.Day,
			"status":    statRes.Group.Status,
		}
		_, err = s.FindOneStat(findCond)
		if err != nil {
			id, err = GCommonMgr.GetMaxId(conf.GConfig.Models.TaskGroup)

			stat = &TaskLogStatMod{
				Id:       id,
				Status:   statRes.Group.Status,
				PlanTime: statRes.Group.Day,
				Count:    int64(statRes.Count),
			}

			s.AddStat(stat)

			continue
		}

		uptData = bson.M{
			"$set": bson.M{
				"status":    statRes.Group.Status,
				"plan_time": statRes.Group.Day,
				"count":     int64(statRes.Count),
			},
		}
		s.UpdateOne(findCond, uptData)
	}
}

func (s *TaskLogStatMgr) LogStat() (res []StatRes, err error) {
	var (
		cur    *mongo.Cursor
		client *mongo.Client
		p      interface{}
	)

	today := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	pipeline := `
		[
		  {
			"$match": {
			  "plan_time": { "$gt": "` + today + `" }
			}
		  },
		  {
			"$project": {
			  "status": 1,
				"count": 1,
			  "day": { "$substr": [ "$plan_time", 0, 10 ] }
			}
		  },
		  {
			"$group": {
			  "_id": { "day": "$day", "status": "$status" },
			  "count": { "$sum": "$count" }
			}
		  },
			{"$sort":{"_id.day":1}}
		]
		`

	opts := options.Aggregate()
	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection := client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLogStat)
	if cur, err = collection.Aggregate(context.TODO(), mdb.MongoPipeline(pipeline), opts); err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &res); err != nil {
		return nil, err
	}

	return
}

/*
添加或者更新文档
*/
func (s *TaskLogStatMgr) UpsertDoc(stat *TaskLogStatMod) (err error) {
	var (
		id       int64
		uptData  bson.M
		findCond bson.M
	)

	findCond = bson.M{
		"plan_time": stat.PlanTime,
		"status":    stat.Status,
	}
	_, err = s.FindOneStat(findCond)

	if err != nil {
		id, err = GCommonMgr.GetMaxId(conf.GConfig.Models.TaskGroup)

		stat.Id = id

		return s.AddStat(stat)
	}

	uptData = bson.M{
		"$inc": bson.M{
			"count": stat.Count,
		},
	}
	return s.UpdateOne(findCond, uptData)
}

func (s *TaskLogStatMgr) UpdateOne(uptCond interface{}, update interface{}) (err error) {
	var (
		collection *mongo.Collection
		client     *mongo.Client
		p          interface{}
		result     *mongo.UpdateResult
		opts       *options.UpdateOptions
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLogStat)

	// 执行删除
	opts = options.Update()
	if result, err = collection.UpdateOne(context.TODO(), uptCond, update, opts); err != nil {
		log.Errorln(err)
		return
	}

	log.Debugln("更新"+conf.GConfig.Models.TaskLogStat+"结果：", result)

	return
}

func (s *TaskLogStatMgr) AddStat(stat *TaskLogStatMod) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertOneResult
		client     *mongo.Client
		p          interface{}
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLogStat)

	if result, err = collection.InsertOne(context.TODO(), stat); err != nil {
		log.Errorln(err)
		return
	}

	log.Debugln(" 插入的 "+conf.GConfig.Models.TaskLogStat+" ID：", result)
	return
}

func (s *TaskLogStatMgr) FindOneStat(findCond interface{}) (stat *TaskLogStatMod, err error) {
	var (
		collection  *mongo.Collection
		res         *mongo.SingleResult
		findOptions *options.FindOneOptions
		findStat    TaskLogStatMod
		client      *mongo.Client
		p           interface{}
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.TaskLogStat)

	findOptions = options.FindOne()
	res = collection.FindOne(context.TODO(), findCond, findOptions)

	if err = res.Decode(&findStat); err != nil {
		return nil, err
	}

	stat = &findStat
	return
}
