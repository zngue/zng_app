package where

import (
	"fmt"
	"reflect"
)

type Operation string

const (
	Eq    Operation = "eq"
	Gt    Operation = "gt"
	Lt    Operation = "lt"
	Gte   Operation = "gte"
	Lte   Operation = "lte"
	Neq   Operation = "neq"
	In    Operation = "in"
	Like  Operation = "like"
	NotIn Operation = "no"
)

// 操作符 转化成符号
func (o Operation) String() string {
	switch o {
	case Eq:
		return "="
	case Gt:
		return ">"
	case Lt:
		return "<"
	case Gte:
		return ">="
	case Lte:
		return "<="
	case Neq:
		return "!="
	case In:
		return "IN"
	case Like:
		return "LIKE"
	case NotIn:
		return "NOT IN"
	default:
		return ""
	}
}

// 常用操作符
var (
	EqOperation    = Eq
	GtOperation    = Gt
	LtOperation    = Lt
	GteOperation   = Gte
	LteOperation   = Lte
	NeqOperation   = Neq
	InOperation    = In
	LikeOperation  = Like
	NotInOperation = NotIn
)

type Option struct {
	Field     string
	Operation Operation
	Value     any
}
type Fn func(opt *Option)

func Where(v any) []*Option {
	var whereOptions []*Option
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
		whereOptions = append(whereOptions, &Option{
			Field:     fileName,
			Operation: Operation(operation),
			Value:     valueInterface,
		})
	}
	return whereOptions
}

func DataWhereOption(filed string, operation Operation, value any) Fn {
	return func(opt *Option) {
		opt.Field = filed
		opt.Operation = operation
		opt.Value = value
	}
}

func NewWhere(opts ...*Option) (where map[string]any) {
	where = make(map[string]any)
	for _, opt := range opts {
		if opt.Value == "" || opt.Value == nil || opt.Value == 0 {
			continue
		}
		if opt.Operation.String() == "" {
			continue
		}
		if opt.Operation == Like {
			where[fmt.Sprintf("%s %s ?", opt.Field, opt.Operation.String())] = fmt.Sprintf("%%%s%%", opt.Value)
		} else {
			where[fmt.Sprintf("%s %s ?", opt.Field, opt.Operation.String())] = opt.Value
		}
	}
	return
}

func NewWhereFn(fns ...Fn) map[string]any {
	var where = make(map[string]any)
	for _, fn := range fns {
		var opt = new(Option)
		fn(opt)
		if opt.Value != nil {
			if opt.Operation == Like {
				where[fmt.Sprintf("%s %s ?", opt.Field, opt.Operation)] = fmt.Sprintf("%%%s%%", opt.Value)
			} else {
				where[fmt.Sprintf("%s %s ?", opt.Field, opt.Operation)] = opt.Value
			}
		}
	}
	return where
}
