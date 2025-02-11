package util

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// StructToMap 将结构体转换为 map[string]string
func StructToMap(v interface{}) map[string]string {
	result := make(map[string]string)
	val := reflect.ValueOf(v)
	// 如果是指针，获取指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		fieldName := fieldType.Name
		// 构建字段名，考虑嵌套结构体
		// 根据不同类型的字段进行处理
		switch field.Kind() {
		case reflect.String:
			result[fieldName] = field.String()
		case reflect.Int:
			result[fieldName] = strconv.Itoa(int(field.Int()))
		case reflect.Float64:
			result[fieldName] = fmt.Sprintf("%.2f", field.Float())
		case reflect.Slice:
			// 对切片类型做特殊处理
			if field.Len() > 0 {
				// 这里假设切片中元素是字符串或可以转换为字符串
				// 不使用带索引的键，而是将每个元素都添加为相同的键名
				for j := 0; j < field.Len(); j++ {
					result[fieldName] = fmt.Sprintf("%v", field.Index(j).Interface())
				}
			}
		case reflect.Struct:
			// 递归处理嵌套结构体
			if field.Type() == reflect.TypeOf(time.Time{}) {
				// 如果是时间类型，使用自定义格式
				result[fieldName] = field.Interface().(time.Time).Format(time.RFC3339)
			} else {
				// 对嵌套结构体递归处理
				subMap := StructToMap(field.Interface())
				for k, vals := range subMap {
					result[k] = vals
				}
			}
		default:
			// 对其他类型字段的通用处理
			result[fieldName] = fmt.Sprintf("%v", field.Interface())
		}
	}

	return result
}
