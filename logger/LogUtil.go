package logger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

type LogUtil struct {
}

// formatMessage 智能格式化消息，支持多种格式
func formatMessage(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	// 只有一个参数
	if len(args) == 1 {
		return valueToString(args[0])
	}

	// 第一个参数是 string
	if format, ok := args[0].(string); ok {

		// Python 风格 {}
		if strings.Contains(format, "{}") {
			return formatPythonStyle(format, args[1:]...)
		}

		// fmt 风格 %
		if strings.Contains(format, "%") {
			return fmt.Sprintf(format, args[1:]...)
		}

		// 普通拼接
		var b strings.Builder
		b.WriteString(format)

		for _, arg := range args[1:] {
			b.WriteByte(' ')
			b.WriteString(valueToString(arg))
		}
		return b.String()
	}

	// 第一个不是 string
	var b strings.Builder
	for i, arg := range args {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(valueToString(arg))
	}
	return b.String()
}

// formatPythonStyle 实现Python风格的{}占位符格式化
func formatPythonStyle(format string, args ...any) string {
	result := format

	for _, arg := range args {
		if !strings.Contains(result, "{}") {
			break
		}
		result = strings.Replace(result, "{}", valueToString(arg), 1)
	}

	return result
}

func valueToString(v any) string {
	if v == nil {
		return "null"
	}

	rv := reflect.ValueOf(v)

	// 解引用指针
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "null"
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct,
		reflect.Map,
		reflect.Slice,
		reflect.Array:
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
	}

	return fmt.Sprint(v)
}

func Info(args ...interface{}) {
	slog.Info(formatMessage(args...))
}

func Warn(args ...interface{}) {
	slog.Warn(formatMessage(args...))
}

func Error(args ...interface{}) {
	slog.Error(formatMessage(args...))
}

func Debug(args ...interface{}) {
	slog.Debug(formatMessage(args...))
}
