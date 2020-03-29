package api

import (
	"cycron/conf"
	"cycron/libs"
	"cycron/mod"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func doSaveUser(resp http.ResponseWriter, req *http.Request) {
	var (
		err       error
		postStr   string
		newPasswd string
		bytes     []byte
		user      mod.UserMod
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "content-type")

	fmt.Println("coming a ", req.Method, " request: ", time.Now())
	if "OPTIONS" == req.Method {
		err = ServerError("忽略 OPTIONS 请求")
		goto ERR
	}

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
		err = ServerError("参数错误")
		goto ERR
	}

	if "" != newPasswd {
		user.Password = libs.Md5(newPasswd)
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
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func doLogin(resp http.ResponseWriter, req *http.Request) {
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

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "content-type")

	fmt.Println("coming a ", req.Method, " request: ", time.Now())
	if "OPTIONS" == req.Method {
		err = ServerError("忽略 OPTIONS 请求")
		goto ERR
	}

	// 1, 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 2, 取表单中的task字段
	email = req.PostForm.Get("email")
	passwd = req.PostForm.Get("password")

	// 4, 保存到mongoDB

	if "" == email || "" == passwd {
		err = ServerError("参数错误")
		goto ERR
	}

	findCond = bson.M{
		"email":    email,
		"password": libs.Md5(passwd),
	}
	if user, err = mod.GUserMgr.FindOneUser(findCond); err != nil {
		goto ERR
	}

	fmt.Println(user)

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
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func userInfo(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		bytes    []byte
		user     *mod.UserMod
		findCond bson.M
		resUser  map[string]interface{}
		token    []string
		ok       bool
		claims   *libs.JwtCustomClaims
	)

	// Stop here if its Preflighted OPTIONS request
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "content-type, token")

	fmt.Println("coming a ", req.Method, " request: ", time.Now())
	if "OPTIONS" == req.Method {
		err = ServerError("忽略 OPTIONS 请求")
		goto ERR
	}

	if token, ok = req.Header["Token"]; !ok {
		err = ServerError("无 token 参数")
		goto ERR_TOKEN
	}

	if claims, err = libs.ParseToken(token[0], []byte(conf.GConfig.Server.JwtSecret)); err != nil {
		goto ERR_TOKEN
	}

	findCond = bson.M{
		"email": claims.Email,
		"_id":   claims.Id,
	}
	if user, err = mod.GUserMgr.FindOneUser(findCond); err != nil {
		goto ERR
	}

	fmt.Println(user)

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
ERR:
	// 6, 返回异常应答
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1001, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
	return
ERR_TOKEN:
	fmt.Println(time.Now(), err)
	if bytes, err = libs.BuildResponse(1000, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
	return
}
