package core

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/zngue/zng_app/core/errors_ez"
)

func EncodeQuery(req any) (url.Values, error) {
	values := make(url.Values)
	if req == nil {
		return values, nil
	}
	value := reflect.ValueOf(req)
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return values, nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return values, nil
	}
	if err := encodeStructQuery(value, values); err != nil {
		return nil, err
	}
	return values, nil
}

func encodeStructQuery(value reflect.Value, values url.Values) error {
	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldValue := value.Field(i)

		name, omitempty := parseJSONTag(field)
		if name == "" {
			continue
		}

		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}
		inner := fieldValue
		for inner.Kind() == reflect.Ptr {
			inner = inner.Elem()
		}

		switch inner.Kind() {
		case reflect.Struct:
			if err := encodeStructQuery(inner, values); err != nil {
				return err
			}
		case reflect.Slice, reflect.Array:
			if omitempty && inner.Len() == 0 {
				continue
			}
			if err := encodeSlice(name, inner, values); err != nil {
				return err
			}
		default:
			if omitempty && isEmptyValue(inner) {
				continue
			}
			if err := encodeScalar(name, inner, values); err != nil {
				return err
			}
		}
	}
	return nil
}

func encodeSlice(name string, value reflect.Value, values url.Values) error {
	if value.Type().Elem().Kind() == reflect.Uint8 {
		values.Add(name, string(value.Bytes()))
		return nil
	}
	for i := 0; i < value.Len(); i++ {
		values.Add(name, fmt.Sprint(value.Index(i).Interface()))
	}
	return nil
}

func encodeScalar(name string, value reflect.Value, values url.Values) error {
	if value.Kind() == reflect.Interface && value.IsNil() {
		return nil
	}
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return nil
	}
	values.Add(name, fmt.Sprint(value.Interface()))
	return nil
}

func parseJSONTag(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false
	}
	parts := strings.Split(tag, ",")
	name := strings.TrimSpace(parts[0])
	if name == "" {
		name = strings.ToLower(field.Name)
	}
	omitempty := false
	for _, option := range parts[1:] {
		if strings.TrimSpace(option) == "omitempty" {
			omitempty = true
			break
		}
	}
	return name, omitempty
}

func isEmptyValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	}
	return false
}

func BindQuery(ctx TransportContext, req any) error {
	if req == nil {
		return nil
	}
	value := reflect.ValueOf(req)
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}
	return bindStructQuery(ctx, value)
}

func bindStructQuery(ctx TransportContext, value reflect.Value) error {
	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		if !field.IsExported() {
			continue
		}
		name, _ := parseJSONTag(field)
		if name == "" {
			continue
		}
		fieldValue := value.Field(i)

		inner := fieldValue
		needsAlloc := false
		for inner.Kind() == reflect.Ptr {
			if inner.IsNil() {
				needsAlloc = true
				break
			}
			inner = inner.Elem()
		}

		if needsAlloc {
			inner = reflect.New(field.Type.Elem()).Elem()
		}

		switch inner.Kind() {
		case reflect.Struct:
			if needsAlloc {
				fieldValue.Set(inner.Addr())
				inner = fieldValue.Elem()
			}
			if err := bindStructQuery(ctx, inner); err != nil {
				return err
			}
		case reflect.Slice, reflect.Array:
			items := ctx.QueryArray(name)
			if len(items) > 0 {
				if err := bindSlice(inner, items); err != nil {
					return errors_ez.WrapF(err, "bind query %s", name)
				}
				if needsAlloc {
					fieldValue.Set(inner.Addr())
				}
			}
		default:
			raw := ctx.Query(name)
			if raw == "" {
				continue
			}
			if err := setFieldValue(inner, raw); err != nil {
				return errors_ez.WrapF(err, "bind query %s", name)
			}
			if needsAlloc {
				fieldValue.Set(inner.Addr())
			}
		}
	}
	return nil
}

func bindSlice(field reflect.Value, items []string) error {
	sliceType := field.Type()
	elementKind := sliceType.Elem().Kind()

	elementIsPtr := elementKind == reflect.Ptr
	elementType := sliceType.Elem()
	if elementIsPtr {
		elementType = elementType.Elem()
	}

	result := reflect.MakeSlice(reflect.SliceOf(sliceType.Elem()), 0, len(items))

	for _, item := range items {
		element := reflect.New(elementType).Elem()
		if err := setFieldValue(element, item); err != nil {
			return err
		}
		if elementIsPtr {
			result = reflect.Append(result, element.Addr())
		} else {
			result = reflect.Append(result, element)
		}
	}

	if result.Len() > 0 || !field.CanSet() {
		if field.IsNil() || result.Len() > 0 {
			field.Set(result)
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, raw string) error {
	if !field.CanSet() {
		return nil
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		field.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(parsed)
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		field.SetFloat(parsed)
	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setFieldValue(field.Elem(), raw)
	default:
		field.SetString(raw)
	}
	return nil
}
