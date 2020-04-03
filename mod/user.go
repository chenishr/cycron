package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserMod struct {
	Id            int64  `bson:"_id"`
	UserName      string `bson:"user_name"`
	Password      string `bson:"password"`
	Email         string `bson:"email"`
	LastLoginTime int64  `bson:"last_login_time"`
	LastIp        string `bson:"last_ip"`
	Status        int    `bson:"status"` // 0 停用；1 正常
	Role          int    `bson:"role"`   // 0 管理；1 普通
	CreateTime    string `bson:"create_time"`
	UpdateTime    string `bson:"update_time"`
}

type UserMgr struct {
}

var (
	GUserMgr *UserMgr
)

func init() {
	GUserMgr = &UserMgr{}
}

/*
添加或者更新文档
*/
func (u *UserMgr) UpsertDoc(user *UserMod) (err error) {
	var (
		id      int64
		uptCond bson.M
		uptData bson.M
	)
	if user.Id == 0 {
		id, err = GCommonMgr.GetMaxId(conf.GConfig.Models.User)

		user.Id = id
		user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

		return u.AddUser(user)
	}

	uptCond = bson.M{"_id": user.Id}
	if "" == user.Password {
		uptData = bson.M{
			"$set": bson.M{
				"user_name":   user.UserName,
				"email":       user.Email,
				"status":      user.Status,
				"role":        user.Role,
				"update_time": time.Now().Format("2006-01-02 15:04:05"),
			},
		}
	} else {
		uptData = bson.M{
			"$set": bson.M{
				"user_name":   user.UserName,
				"email":       user.Email,
				"password":    user.Password,
				"status":      user.Status,
				"role":        user.Role,
				"update_time": time.Now().Format("2006-01-02 15:04:05"),
			},
		}
	}
	return u.UpdateOne(uptCond, uptData)
}

func (u *UserMgr) UpdateOne(uptCond interface{}, update interface{}) (err error) {
	var (
		res        *mongo.UpdateResult
		collection *mongo.Collection
		client     *mongo.Client
		p          interface{}
	)

	log.Traceln("执行更新")

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.User)

	// 执行删除
	if res, err = collection.UpdateOne(context.TODO(), uptCond, update); err != nil {
		log.Errorln(err)
		return
	}

	log.Debugln(res)

	return
}

func (u *UserMgr) AddUser(user *UserMod) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertOneResult
		client     *mongo.Client
		p          interface{}
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.User)

	if result, err = collection.InsertOne(context.TODO(), user); err != nil {
		log.Errorln(err)
		return
	}

	log.Debugln(" 插入的"+conf.GConfig.Models.User+"ID：", result)

	return
}

func (u *UserMgr) FindOneUser(findCond interface{}) (user *UserMod, err error) {
	var (
		collection  *mongo.Collection
		res         *mongo.SingleResult
		findOptions *options.FindOneOptions
		userData    UserMod
		client      *mongo.Client
		p           interface{}
	)

	p, err = dbs.GMongoPool.Get()
	defer dbs.GMongoPool.Put(p)

	client = p.(*mongo.Client)
	collection = client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.User)

	findOptions = options.FindOne()
	res = collection.FindOne(context.TODO(), findCond, findOptions)

	if err = res.Decode(&userData); err != nil {
		return nil, err
	}

	user = &userData
	return
}
