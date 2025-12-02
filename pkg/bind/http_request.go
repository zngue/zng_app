package bind

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	"github.com/go-resty/resty/v2"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zngue/zng_app/errors"
)

type ClientServer interface {
	POST(ctx context.Context, path string, in any, v any) error
	GET(ctx context.Context, path string, in any, v any) error
}
type Client struct {
	BaseURL string
	Version string
	Header  map[string]string // 请求头
}
type ClientOption struct {
	ServerName string            // 服务名称
	Port       string            // 服务端口
	Version    string            // 服务版本
	Header     map[string]string // 请求头
	IsNacos    bool              // 是否使用nacos
	GroupName  string            // nacos分组
	BaseUrl    string
}
type ClientOptionFn func(rs *ClientOption)

func DataWithHeaders(headers map[string]string) ClientOptionFn {
	return func(response *ClientOption) {
		response.Header = headers
	}
}
func DataWithSeviceName(name string) ClientOptionFn {
	return func(response *ClientOption) {
		response.ServerName = name
	}
}
func DataWithPort(port string) ClientOptionFn {
	return func(response *ClientOption) {
		response.Port = port
	}
}
func DataWithVersion(version string) ClientOptionFn {
	return func(response *ClientOption) {
		response.Version = version
	}
}
func DataWithGroupName(name string) ClientOptionFn {
	return func(response *ClientOption) {
		response.GroupName = name
	}
}
func DataWithBaseUrl(url string) ClientOptionFn {
	return func(response *ClientOption) {
		response.BaseUrl = url
	}
}
func DataWithNacosServer(srv naming_client.INamingClient, serviceName, groupName string) ClientOptionFn {
	return func(response *ClientOption) {
		response.IsNacos = true
		var baseUrl, err = NacosServer(srv, serviceName, groupName)
		if err != nil {
			panic(err)
		}
		response.BaseUrl = baseUrl
	}
}

// Server naming_client.INamingClient
func NacosServer(srv naming_client.INamingClient, serviceName, groupName string) (baseURL string, err error) {
	var info model.Service
	info, err = srv.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
		GroupName:   groupName,
	})
	if err != nil {
		return
	}
	if len(info.Hosts) > 0 {
		for _, v := range info.Hosts {
			baseURL = fmt.Sprintf("http://%s:%d", serviceName, v.Port)
			return
		}
	}
	return
}
func NewClientServer(option *ClientOption, fns ...ClientOptionFn) (cli ClientServer, err error) {
	if option == nil {
		option = &ClientOption{
			Version: "v2",
		}
	}
	for _, fn := range fns {
		fn(option)
	}
	cli = &Client{
		BaseURL: option.BaseUrl,
		Version: option.Version,
		Header:  option.Header,
	}
	if option.BaseUrl == "" {
		err = errors.New("nacos server url is empty")
		return
	}
	return
}

func NewNacosClientServer(srv naming_client.INamingClient, serviceName, groupName string, fns ...ClientOptionFn) (cli ClientServer, err error) {
	var option = &ClientOption{
		Version: "v1",
	}
	option.BaseUrl, err = NacosServer(srv, serviceName, groupName)
	if err != nil {
		return
	}
	for _, fn := range fns {
		fn(option)
	}
	cli = &Client{
		BaseURL: option.BaseUrl,
		Version: option.Version,
		Header:  option.Header,
	}
	if option.BaseUrl == "" {
		err = errors.New("nacos server url is empty")
		return
	}
	return
}

// ----------------------
//
//	POST JSON
//
// ----------------------
func (c *Client) POST(ctx context.Context, path string, in any, v any) (err error) {
	client := resty.New()
	request := client.R().SetBody(in)
	request = request.SetContext(ctx)
	var rs *resty.Response
	var url = fmt.Sprintf("%s/%s/%s", c.BaseURL, c.Version, path)
	rs, err = request.Post(url)
	if err != nil {
		return
	}
	if rs.StatusCode() != 200 {
		err = fmt.Errorf("status code: %d", rs.StatusCode())
		return
	}
	s := rs.Body()
	err = json.Unmarshal(s, v)
	if err != nil {
		return
	}
	return nil
}

// GET 请求方法
func (c *Client) GET(ctx context.Context, path string, query any, v any) (err error) {
	params := StructToURLValues(query)
	client := resty.New()
	request := client.R()
	request = request.SetContext(ctx)
	if len(params) > 0 {
		request = request.SetQueryParamsFromValues(params)
	}
	var url = fmt.Sprintf("%s/%s/%s", c.BaseURL, c.Version, path)
	var rs *resty.Response
	rs, err = request.Get(url)
	if err != nil {
		return
	}
	if rs.StatusCode() != 200 {
		err = fmt.Errorf("status code: %d", rs.StatusCode())
		return
	}
	s := rs.Body()
	err = json.Unmarshal(s, v)
	if err != nil {
		return
	}
	return
}
func StructToURLValues(v any) url.Values {
	values := url.Values{}
	parseToURLValues(reflect.ValueOf(v), values)
	return values
}

func parseToURLValues(val reflect.Value, out url.Values) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fType := t.Field(i)

		// ---- 取 json tag ----
		jsonTag := fType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		key := jsonTag
		// 匿名字段：直接展开，不使用 tag
		if fType.Anonymous {
			parseToURLValues(field, out)
			continue
		}
		switch field.Kind() {
		case reflect.Struct:
			// 嵌套结构体提升到顶层
			parseToURLValues(field, out)
		case reflect.Slice, reflect.Array:
			for j := 0; j < field.Len(); j++ {
				out.Add(key, fmt.Sprintf("%v", field.Index(j).Interface()))
			}
		default:
			out.Add(key, fmt.Sprintf("%v", field.Interface()))
		}
	}
}
