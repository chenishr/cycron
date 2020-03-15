package mod

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserMod struct {
	Id            primitive.ObjectID `bson:"_id"`
	UserName      string             `bson:"user_name"`
	Password      string             `bson:"password"`
	Email         string             `bson:"email"`
	LastLoginTime int64              `bson:"last_login_time"`
	LastIp        string             `bson:"last_ip"`
	Status        int                `bson:"status"` // 0 停用；1 正常
	CreateTime    int64             `bson:"create_time"`
	UpdateTime    int64             `bson:"update_time"`
}
