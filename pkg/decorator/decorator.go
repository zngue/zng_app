package decorator

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func NewDecorated[T any](impl T) (rs T) {
	rs, _ = Decorate(impl)
	return
}

// WrapError 包装 error，包含接口名、方法名、调用行号
func WrapError(component, method string, err error) error {
	if err == nil {
		return nil
	}

	// 获取调用堆栈
	pc, file, line, ok := runtime.Caller(2) // 2 层上
	funcName := "unknown"
	if ok {
		f := runtime.FuncForPC(pc)
		if f != nil {
			funcName = f.Name()
			// 去掉路径，只保留函数名
			parts := strings.Split(funcName, ".")
			funcName = parts[len(parts)-1]
		}
		file = trimFilePath(file)
	}

	return fmt.Errorf("[%s.%s.%s] %s:%d %v", component, method, funcName, file, line, err)
}

func trimFilePath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return path
}
func Decorate[T any](impl T) (rs T, err error) {
	var ifacePtr T
	ifaceType := reflect.TypeOf(&ifacePtr).Elem() // interface type
	implValue := reflect.ValueOf(impl)
	if ifaceType.Kind() != reflect.Interface {
		err = errors.New("Decorate[T]: T must be an interface")
		return
	}
	// 类型断言检查
	if !reflect.TypeOf(impl).Implements(ifaceType) {
		err = fmt.Errorf("Decorate[T]: impl type %s does not implement interface %s", reflect.TypeOf(impl).String(), ifaceType.String())
		return
	}
	componentName := ifaceType.Name()
	proxy := createProxy(ifaceType, implValue, componentName)
	rs = proxy.Interface().(T)
	return
}

func createProxy(ifaceType reflect.Type, implValue reflect.Value, name string) reflect.Value {
	numMethods := ifaceType.NumMethod()
	methods := make([]reflect.StructField, numMethods)
	for i := 0; i < numMethods; i++ {
		m := ifaceType.Method(i)
		methods[i] = reflect.StructField{
			Name: m.Name,
			Type: m.Type,
		}
	}
	proxyType := reflect.StructOf(methods)
	proxyValue := reflect.New(proxyType).Elem()

	for i := 0; i < numMethods; i++ {
		m := ifaceType.Method(i)
		methodName := m.Name
		implMethod := implValue.MethodByName(methodName)
		if !implMethod.IsValid() {
			panic(fmt.Errorf("impl missing method: %s", methodName))
		}

		fn := reflect.MakeFunc(m.Type, wrapMethodWithTrace(name, methodName, implMethod))
		proxyValue.Field(i).Set(fn)
	}

	return proxyValue.Convert(ifaceType)
}
func wrapMethodWithTrace(component, method string, implMethod reflect.Value) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) (results []reflect.Value) {
		defer func() {
			if r := recover(); r != nil {
				// panic → error
				err := fmt.Errorf("[%s.%s] panic: %v", component, method, r)
				lastIdx := len(results) - 1
				results[lastIdx] = errToReflect(results[lastIdx].Type(), err)
			}
		}()

		results = implMethod.Call(args)

		// 包装 error（假设最后一个返回值是 error）
		lastIdx := len(results) - 1
		if lastIdx >= 0 {
			if err, ok := results[lastIdx].Interface().(error); ok && err != nil {
				results[lastIdx] = errToReflect(results[lastIdx].Type(), WrapError(component, method, err))
			}
		}

		return results
	}
}

// 将 error 转换成 reflect.Value
func errToReflect(errType reflect.Type, err error) reflect.Value {
	if err == nil {
		return reflect.Zero(errType)
	}
	return reflect.ValueOf(err)
}
