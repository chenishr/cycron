package mod

import (
	"context"
	"cycron/conf"
	"cycron/dbs"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMgr struct {

}

var(
	GUserMgr *UserMgr
)

func init() {
	GUserMgr = &UserMgr{}
}

func (u *UserMgr)AddUser(user *UserMod) (err error) {
	var (
		collection *mongo.Collection
		result     *mongo.InsertOneResult
	)

	collection = dbs.GMongo.Client.Database(conf.GConfig.Models.Db).Collection(conf.GConfig.Models.User)

	if result, err = collection.InsertOne(context.TODO(), user); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" 插入的 ID：", result)
	return
}
