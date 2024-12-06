package log

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func RequestGinLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := make(map[string]any)
		var body []byte
		method := ctx.Request.Method
		route := ctx.FullPath()
		ip := ctx.ClientIP()
		var reqData = make(map[string]any)
		switch method {
		case http.MethodGet:
			for key, values := range ctx.Request.URL.Query() {
				if len(values) > 0 {
					params[key] = values[0] // 只取第一个值
				}
			}
		case http.MethodPost:
			if ctx.Request.URL.Query() != nil {
				for key, values := range ctx.Request.URL.Query() {
					if len(values) > 0 {
						params[key] = values[0] // 只取第一个值
					}
				}
			}
			_ = ctx.Request.ParseForm()
			if ctx.Request.PostForm != nil {
				for key, values := range ctx.Request.PostForm {
					if len(values) > 0 {
						params[key] = values[0] // 只取第一个值
					}
				}
			}
			body, _ = io.ReadAll(ctx.Request.Body)
			var jsonBody = make(map[string]any)
			if len(body) > 0 {
				_ = json.Unmarshal(body, &jsonBody)
			}
			if len(jsonBody) > 0 {
				for key, value := range jsonBody {
					reqData[key] = value
				}
			}
		}
		if len(params) == 0 && len(route) == 0 {
			ctx.Next()
			return
		}
		data := log(map[string]any{
			"method": method,
			"route":  route,
			"ip":     ip,
			"params": params,
		})
		Default().Info("app_request", data...)
		ctx.Next()
	}
}
