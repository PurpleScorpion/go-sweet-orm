package mapperDemo

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-orm/v2/mapper"
	"reflect"
	"strings"
	"testing"
)

func TestDemo2(t *testing.T) {
	var list []Logs
	fmt.Println(getTableId(&Logs{}))
	fmt.Println(getTableId(Logs{}))
	fmt.Println(getTableId(getLogList()))
	fmt.Println(getTableId(getLogListNoPoint()))
	fmt.Println(getTableId(&list))
	fmt.Println("-------------------------------------")
	fmt.Println(getTableName(getLogList()))
	fmt.Println(getTableName(getLogListNoPoint()))
	fmt.Println(getTableName(&Logs{}))
	fmt.Println(getTableName(Logs{}))

	fmt.Println(getTableName(&list))
}
func TestDemo1(t *testing.T) {

	wrapper := mapper.BuilderUpdateWrapper(getLogList())
	wrapper.InsertAll(10, true)
}

func getLogList() []*Logs {
	var list []*Logs
	var i int32 = 1

	for ; i < 104; i++ {
		list = append(list, &Logs{LogContent: fmt.Sprintf("test%d", i), LogLevel: fmt.Sprintf("info%d", i)})
	}
	return list
}

func getLogListNoPoint() []Logs {
	var list []Logs
	var i int32 = 1

	for ; i < 10; i++ {
		list = append(list, Logs{LogContent: fmt.Sprintf("test%d", i), LogLevel: fmt.Sprintf("info%d", i)})
	}
	return list
}

func getTableId(T any) string {
	v := reflect.Indirect(reflect.ValueOf(T))
	//加入缓存 , 速度更快
	tp := v.Type().String()

	switch v.Kind() {
	case reflect.Slice:
		return getTableId4Array(T, tp)
	case reflect.Struct:
		return getTableId4Object(T, tp)
	default:
		fmt.Println("Input is neither a slice nor a struct")
	}
	return "null"
}

func getTableId4Array(T any, typeName string) string {
	v := reflect.Indirect(reflect.ValueOf(T))
	elemType := v.Type().Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	for i := 0; i < elemType.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := elemType.Field(i)
		// 通过字段名, 找到字段类型信息
		if tp, ok := elemType.FieldByName(fieldType.Name); ok {
			// 从tag中取出需要的tag
			if !isEmpty(tp.Tag.Get("tableId")) {
				idName := camelToUnderscore(tp.Tag.Get("tableId"))
				return idName
			}
		}
	}
	return "null"
}

func getTableId4Object(T any, typeName string) string {
	ref := reflect.Indirect(reflect.ValueOf(T)).Type()

	for i := 0; i < ref.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := ref.Field(i)
		// 通过字段名, 找到字段类型信息
		if tp, ok := ref.FieldByName(fieldType.Name); ok {
			// 从tag中取出需要的tag
			if !isEmpty(tp.Tag.Get("tableId")) {
				idName := camelToUnderscore(tp.Tag.Get("tableId"))
				return idName
			}
		}
	}
	return "null"
}

func getTableName(T interface{}) string {
	funcName := "TableName"
	tp := getParmaStruct(T)
	if tp == "null" {
		return "null"
	}

	if tp == "Object" {
		methodValue := reflect.ValueOf(T).MethodByName(funcName)
		if !methodValue.IsValid() {
			return "null"
		}
		back := methodValue.Call(nil)
		tbName := back[0].String()
		return tbName
	} else {
		v := reflect.Indirect(reflect.ValueOf(T))

		elemType := v.Type().Elem()
		// elemType := v.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		obj := reflect.New(elemType).Elem().Interface()
		funcVal := reflect.ValueOf(obj).MethodByName(funcName)
		// methodValue, _ := elemType.MethodByName(funcName)
		// fmt.Println(methodValue.Name)
		// funcVal := methodValue.Func

		if !funcVal.IsValid() {
			return "null"
		}
		back := funcVal.Call(nil)
		tbName := back[0].String()
		return tbName
	}

}

func getParmaStruct(T any) string {
	v := reflect.Indirect(reflect.ValueOf(T))
	switch v.Kind() {
	case reflect.Slice:
		return "Array"
	case reflect.Struct:
		return "Object"
	default:
		fmt.Println("Input is neither a slice nor a struct")
	}
	return "null"
}

func camelToUnderscore(str string) string {

	var result strings.Builder
	// 将字符串首字母转换为小写
	str = strings.ToLower(str[0:1]) + str[1:]

	for i, char := range str {
		if char >= 'A' && char <= 'Z' {
			// 如果当前字符是大写字母，则在前面添加下划线，并将该字母转为小写
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteByte(byte(char + 32))
		} else {
			result.WriteByte(byte(char))
		}
	}

	return result.String()
}
