package middleware

import (
	error2 "cycron/api/error"
	"cycron/conf"
	"cycron/libs"
	"cycron/mod"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

var (
	User *mod.UserMod
)

//  设置响应头部
func CheckToken() Middleware {

	// 创建中间件
	return func(f http.HandlerFunc) http.HandlerFunc {

		// 创建一个新的handler包装http.HandlerFunc
		return func(resp http.ResponseWriter, req *http.Request) {
			// 中间件的处理逻辑

			var (
				findCond bson.M
				token    []string
				ok       bool
				claims   *libs.JwtCustomClaims
				err      error
				bytes    []byte
			)

			if token, ok = req.Header["Token"]; !ok {
				err = error2.ServerError("无 token 参数")
				goto ERR_TOKEN
			}

			if claims, err = libs.ParseToken(token[0], []byte(conf.GConfig.Server.JwtSecret)); err != nil {
				goto ERR_TOKEN
			}

			findCond = bson.M{
				"email": claims.Email,
				"_id":   claims.Id,
			}
			if User, err = mod.GUserMgr.FindOneUser(findCond); err != nil {
				goto ERR_TOKEN
			}

			// 正常
			// 调用下一个中间件或者最终的handler处理程序
			f(resp, req)
			return

		ERR_TOKEN:
			// 异常
			log.Errorln(err)
			if bytes, err = libs.BuildResponse(1000, err.Error(), nil); err == nil {
				resp.Write(bytes)
			}
			return
		}
	}
}
