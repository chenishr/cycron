package middleware

import (
	"net/http"
)

//  设置响应头部
func SetHeader() Middleware {

	// 创建中间件
	return func(f http.HandlerFunc) http.HandlerFunc {

		// 创建一个新的handler包装http.HandlerFunc
		return func(resp http.ResponseWriter, req *http.Request) {

			// 中间件的处理逻辑

			resp.Header().Set("Access-Control-Allow-Origin", "*")
			resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			resp.Header().Set("Access-Control-Allow-Headers", "content-type, token")

			// 调用下一个中间件或者最终的handler处理程序
			f(resp, req)
		}
	}
}
