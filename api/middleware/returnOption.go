package middleware

import (
	"cycron/libs"
	"net/http"
)

//  设置响应头部
func ReturnOption() Middleware {

	// 创建中间件
	return func(f http.HandlerFunc) http.HandlerFunc {

		// 创建一个新的handler包装http.HandlerFunc
		return func(resp http.ResponseWriter, req *http.Request) {

			// 中间件的处理逻辑
			if "OPTIONS" == req.Method {
				if bytes, err := libs.BuildResponse(1001, "忽略 OPTIONS 请求", nil); err == nil {
					resp.Write(bytes)
				}
			} else {
				// 调用下一个中间件或者最终的handler处理程序
				f(resp, req)
			}
		}
	}
}
