package bind

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type User struct {
	Name     string         `json:"name"`
	Age      int            `json:"age"`
	UserInfo *UserInfo      `json:"userInfo"`
	Mapping  map[string]any `json:"mapping"`
}
type UserInfo struct {
	Email     string   `json:"email"`
	EmailName []string `json:"emailName"`
	Points    []int    `json:"points"`
}

func Test_HttpGet(t *testing.T) {
	var data = &User{
		Name: "333",
		Age:  25,
		Mapping: map[string]any{
			"name5": map[string]any{
				"userName5": 7,
				"userName6": 8,
			},
		},
		UserInfo: &UserInfo{
			Email: "sds",
			EmailName: []string{
				"sds",
				"sds_1",
			},
			Points: []int{
				1,
				2,
				3,
			},
		},
	}
	var dataRs = StructToURLValuesIn(data)
	fmt.Println(dataRs)
}

// 主入口函数 - 扁平化输出，不添加前缀
func StructToURLValuesIn(obj any) url.Values {
	values := url.Values{}
	encodeStruct(obj, &values)
	return values
}

// 递归编码函数
func encodeStruct(obj any, values *url.Values) {
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
		// 获取 JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue // 跳过没有 json tag 或标记为忽略的字段
		}
		// 处理 omitempty 和其他选项
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
