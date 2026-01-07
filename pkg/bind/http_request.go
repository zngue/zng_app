package bind

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

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
	BaseURL       string
	Version       string
	Header        map[string]string // 请求头
	ServiceName   string
	Authorization string
}
type ClientOption struct {
	ServerName    string            // 服务名称
	Port          string            // 服务端口
	Version       string            // 服务版本
	Header        map[string]string // 请求头
	IsNacos       bool              // 是否使用nacos
	GroupName     string            // nacos分组
	BaseUrl       string
	Authorization string
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
		option = &ClientOption{}
	}
	for _, fn := range fns {
		fn(option)
	}
	cli = &Client{
		BaseURL:       option.BaseUrl,
		Version:       option.Version,
		Header:        option.Header,
		ServiceName:   option.ServerName,
		Authorization: option.Authorization,
	}
	if option.BaseUrl == "" {
		err = errors.New("BaseUrl url is empty")
		return
	}
	return
}

func NewNacosClientServer(srv naming_client.INamingClient, serviceName, groupName string, fns ...ClientOptionFn) (cli ClientServer, err error) {
	var option = &ClientOption{}
	option.BaseUrl, err = NacosServer(srv, serviceName, groupName)
	if err != nil {
		return
	}
	for _, fn := range fns {
		fn(option)
	}
	cli = &Client{
		BaseURL:       option.BaseUrl,
		Version:       option.Version,
		Header:        option.Header,
		ServiceName:   option.ServerName,
		Authorization: option.Authorization,
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
	var request = c.RequestCommon(ctx)
	var rs *resty.Response
	var url = fmt.Sprintf("%s%s", c.BaseURL, path)
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

func (c *Client) RequestCommon(ctx context.Context) *resty.Request {
	client := resty.New()
	request := client.R()
	request = request.SetContext(ctx)
	var udid = FromUDIDContext(ctx)
	if udid == "" {
		udid = OriginUDID()
	}
	request = request.SetHeader(RequestFromService, FromServerLocalContext(ctx))
	if c.Version != "" {
		request = request.SetHeader(RequestVersion, c.Version)
	}
	if c.Authorization != "" {
		request = request.SetHeader(RequestAuthorization, c.Authorization)
	}
	request = request.SetHeader(RequestIDKey, udid)
	return request
}

// GET 请求方法
func (c *Client) GET(ctx context.Context, path string, query any, v any) (err error) {
	params := StructToURLValues(query)
	var request = c.RequestCommon(ctx)
	if len(params) > 0 {
		request = request.SetQueryParamsFromValues(params)
	}
	var url = fmt.Sprintf("%s%s", c.BaseURL, path)
	var rs *resty.Response
	request.URL = url
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
	EncodeStruct(v, &values)
	return values
}

func EncodeStruct(obj any, values *url.Values) {
	v := reflect.ValueOf(obj)
	// 处理指针
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	// 必须是结构体
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonKey := strings.Split(jsonTag, ",")[0]
		if jsonKey == "" {
			continue
		}
		// 处理指针
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}
		ValuesData(jsonKey, fieldValue, values)
	}
}
func ValuesData(key string, fieldValue reflect.Value, values *url.Values) {
	switch fieldValue.Kind() {
	case reflect.String:
		str := fieldValue.String()
		if str != "" {
			values.Add(key, str)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		values.Add(key, strconv.FormatInt(fieldValue.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		values.Add(key, strconv.FormatUint(fieldValue.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		values.Add(key, strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64))
	case reflect.Bool:
		values.Add(key, strconv.FormatBool(fieldValue.Bool()))
	case reflect.Slice, reflect.Array:
		if fieldValue.Len() == 0 {
			return
		}
		for i := 0; i < fieldValue.Len(); i++ {
			elem := fieldValue.Index(i)
			if elem.Kind() == reflect.Ptr {
				if elem.IsNil() {
					continue
				}
				elem = elem.Elem()
			}
			ValuesData(key, elem, values)
		}
	case reflect.Map:
		if fieldValue.Len() == 0 {
			return
		}
		for _, keyValue := range fieldValue.MapKeys() {
			mapValue := fieldValue.MapIndex(keyValue)
			var mapKey = keyValue.String()
			ValuesData(mapKey, mapValue, values)
		}
	case reflect.Struct: // 递归处理嵌套结构体，但不添加前缀
		EncodeStruct(fieldValue.Interface(), values)
	case reflect.Interface: // 处理 interface{} 类型
		if !fieldValue.IsNil() {
			v := reflect.ValueOf(fieldValue.Interface())
			ValuesData(key, v, values)
		}
	}
}
