package handle

import (
	error2 "cycron/api/error"
	"cycron/api/middleware"
	"cycron/conf"
	"cycron/libs"
	"cycron/mod"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func DoSaveUser(resp http.ResponseWriter, req *http.Request) {
	var (
		err       error
		postStr   string
		newPasswd string
		bytes     []byte
		user      mod.UserMod
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	postStr = req.PostForm.Get("user")
	newPasswd = req.PostForm.Get("newPassword")

	// 3, 反序列化task
	if err = json.Unmarshal([]byte(postStr), &user); err != nil {
		goto ERR
	}

	// 4, 保存到mongoDB

	if "" == user.UserName || "" == user.Email {
		err = error2.ServerError("参数错误")
		goto ERR
	}

	if "" != newPasswd {
		user.Password = libs.Md5(newPasswd)
	} else {
		user.Password = ""
	}

	if err = mod.GUserMgr.UpsertDoc(&user); err != nil {
		goto ERR
	}

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func DoLogin(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		email    string
		passwd   string
		bytes    []byte
		user     *mod.UserMod
		findCond bson.M
		jwtToken string
		resData  map[string]interface{}
		resUser  map[string]interface{}
	)

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	email = req.PostForm.Get("email")
	passwd = req.PostForm.Get("password")

	// 4, 保存到mongoDB

	if "" == email {
		err = error2.ServerError("参数错误")
		goto ERR
	}

	findCond = bson.M{
		"email":    email,
		"password": libs.Md5(passwd),
	}
	if user, err = mod.GUserMgr.FindOneUser(findCond); err != nil {
		goto ERR
	}

	if jwtToken, err = libs.CreateToken([]byte(conf.GConfig.Server.JwtSecret), uint(user.Id), user.UserName, user.Email); err != nil {
		goto ERR
	}

	user.Password = ""
	resData = make(map[string]interface{})
	resData["token"] = jwtToken

	resUser = make(map[string]interface{})
	resUser["Id"] = user.Id
	resUser["UserName"] = user.UserName
	resUser["Email"] = user.Email
	resUser["Role"] = user.Role
	resUser["Status"] = user.Status

	resData["user"] = resUser

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", resData); err == nil {
		resp.Write(bytes)
		return
	}
ERR:
	// 6, 返回异常应答
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func UserInfo(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		bytes   []byte
		user    *mod.UserMod
		resUser map[string]interface{}
	)

	user = middleware.User

	resUser = make(map[string]interface{})
	resUser["Id"] = user.Id
	resUser["UserName"] = user.UserName
	resUser["Email"] = user.Email
	resUser["Role"] = user.Role
	resUser["Status"] = user.Status

	// 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = libs.BuildResponse(0, "success", resUser); err == nil {
		resp.Write(bytes)
		return
	}

	// 6, 返回异常应答
	log.Errorln(err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
	return
}
