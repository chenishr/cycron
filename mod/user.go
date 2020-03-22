package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMod struct {
	Id            int64  `bson:"_id"`
	UserName      string `bson:"user_name"`
	Password      string `bson:"password"`
	Email         string `bson:"email"`
	LastLoginTime int64  `bson:"last_login_time"`
	LastIp        string `bson:"last_ip"`
	Status        int    `bson:"status"` // 0 停用；1 正常
	CreateTime    int64  `bson:"create_time"`
	UpdateTime    int64  `bson:"update_time"`
}

type UserMgr struct {
}

var (
	GUserMgr *UserMgr
)

func init() {
	GUserMgr = &UserMgr{}
}

func (u *UserMgr) AddUser(user *UserMod) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertOneResult
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.User)

	if result, err = collection.InsertOne(context.TODO(), user); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 插入的 user ID：", result)
	return
}
