package bind

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zngue/zng_app/pkg/util"
)

type Http[V any] struct {
	Code int    `json:"code" `
	Msg  string `json:"msg" `
	Data V      `json:"data" `
}

func HttpGetRequest[V any](url, path string, i any) (rl V, err error) {
	params := util.StructToMap(i)
	client := resty.New()
	request := client.R()
	if len(params) > 0 {
		request = request.SetQueryParams(params)
	}
	var rs *resty.Response
	rs, err = request.Get(fmt.Sprintf("%s%s", url, path))
	if err != nil {
		return
	}
	s := rs.Body()
	var val Http[V]
	err = json.Unmarshal(s, &val)
	if err != nil {
		return
	}
	if val.Code != 200 {
		err = fmt.Errorf("请求失败,code:%d,msg:%s", val.Code, val.Msg)
	}
	if err != nil {
		return
	}
	rl = val.Data
	return
}
func HttpPostRequest[V any](url, path string, i any) (rl V, err error) {
	client := resty.New()
	request := client.R().SetBody(i)
	var rs *resty.Response
	rs, err = request.Post(fmt.Sprintf("%s%s", url, path))
	if err != nil {
		return
	}
	s := rs.Body()
	var val Http[V]
	err = json.Unmarshal(s, &val)
	if err != nil {
		return
	}
	if val.Code != 200 {
		err = fmt.Errorf("请求失败,code:%d,msg:%s", val.Code, val.Msg)
	}
	if err != nil {
		return
	}
	rl = val.Data
	return
}
