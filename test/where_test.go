package test

import (
	"github.com/zngue/zng_app/db/data/page"
	"github.com/zngue/zng_app/db/data/where"
	"reflect"
	"testing"
)

type UserRequest struct {
	Name string `form:"name" where:"like" field:"name"`
	Age  int    `form:"age"  where:"eq" field:"age"`
	Id   int    `form:"id" where:"gt" field:"id"`
	Ids  []int  `form:"ids" where:"in" field:"id"`
}

// 定义一个结构体
type Person struct {
	Name string `form:"name" field:"name" where:"eq"`
	Age  int    `field:"age" where:"gt"`
	City string `form:"city" field:"city" where:"like"`
	User string `form:"user" field:"user" where:"like"`
}

func Where(v any) []*WhereOption {
	var whereOptions []*WhereOption
	refType := reflect.ValueOf(v)
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
	}
	for i := 0; i < refType.NumField(); i++ {
		valueType := refType.Type().Field(i)
		valueInterface := refType.Field(i).Interface()
		if valueType.Type.Kind() == reflect.Struct || valueType.Type.Kind() == reflect.Ptr {
			vals := Where(valueInterface)
			whereOptions = append(whereOptions, vals...)
			continue
		}
		fileName := valueType.Tag.Get("field")
		operation := valueType.Tag.Get("where")
		if fileName == "" || operation == "" {
			continue
		}
		whereOptions = append(whereOptions, &WhereOption{
			Filed:     fileName,
			Operation: operation,
			Value:     valueInterface,
		})
	}
	return whereOptions
}

// 反射获取参数
func TestReflect(t *testing.T) {
	p := Person{Name: "Alice", Age: 30, City: "New York", User: "zng"}
	options := where.Where(p)
	newWhere := where.NewWhere(options...)
	t.Log(options, newWhere)
	//var request = &UserRequest{
	//	Name: "zng",
	//}
	//whereOptions := Where(&request)
	//t.Log(whereOptions)
}

type WhereOption struct {
	Filed     string
	Operation string
	Value     any
}

func TestWhere(t *testing.T) {

	pageData := page.NewPage(
		page.DataWithPage(1),
		page.DataWithPageSize(10),
	)
	t.Log(pageData)

}
