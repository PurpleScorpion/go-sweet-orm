package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type JSONObject struct {
	Jm map[string]interface{}
}

type JSONArray struct {
	Jm []interface{}
}

func (js *JSONObject) GetData() map[string]interface{} {
	return js.Jm
}
func (js *JSONObject) ParseObject(str string) *JSONObject {
	var data map[string]interface{}
	json.Unmarshal([]byte(str), &data)
	js.Jm = data
	return js
}

func (js *JSONObject) GetTime(key string) (time.Time, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return time.Time{}, false
	}
	switch t := value.(type) {
	case time.Time:
		val := value.(time.Time)
		return val, true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return time.Time{}, false
}

func (js *JSONObject) GetString(key string) (string, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return "", false
	}
	switch t := value.(type) {
	case string:
		val := value.(string)
		return val, true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return "", false
}

func (js *JSONObject) GetInt(key string) (int, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return 0, false
	}
	switch t := value.(type) {
	case int:
		val := value.(int)
		return val, true
	case int16:
		val := value.(int16)
		return int(val), true
	case int32:
		val := value.(int32)
		return int(val), true
	case int64:
		val := value.(int64)
		return int(val), true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return 0, false
}

func (js *JSONObject) GetInt16(key string) (int16, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return 0, false
	}
	switch t := value.(type) {
	case int:
		val := value.(int)
		return int16(val), true
	case int16:
		val := value.(int16)
		return int16(val), true
	case int32:
		val := value.(int32)
		return int16(val), true
	case int64:
		val := value.(int64)
		return int16(val), true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return 0, false
}

func (js *JSONObject) GetInt32(key string) (int32, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return 0, false
	}
	switch t := value.(type) {
	case int:
		val := value.(int)
		return int32(val), true
	case int16:
		val := value.(int16)
		return int32(val), true
	case int32:
		val := value.(int32)
		return int32(val), true
	case int64:
		val := value.(int64)
		return int32(val), true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return 0, false
}

func (js *JSONObject) GetInt64(key string) (int64, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return 0, false
	}

	switch t := value.(type) {
	case int:
		val := value.(int)
		return int64(val), true
	case int16:
		val := value.(int16)
		return int64(val), true
	case int32:
		val := value.(int32)
		return int64(val), true
	case int64:
		val := value.(int64)
		return int64(val), true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return 0, false
}

func (js *JSONObject) GetFloat32(key string) (float32, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return 0, false
	}
	switch t := value.(type) {
	case int:
		val := value.(int)
		return float32(val), true
	case int16:
		val := value.(int16)
		return float32(val), true
	case int32:
		val := value.(int32)
		return float32(val), true
	case int64:
		val := value.(int64)
		return float32(val), true
	case float32:
		val := value.(float32)
		return float32(val), true
	case float64:
		val := value.(float64)
		return float32(val), true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return 0, false
}

func (js *JSONObject) GetFloat64(key string) (float64, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return 0, false
	}
	switch t := value.(type) {
	case int:
		val := value.(int)
		return float64(val), true
	case int16:
		val := value.(int16)
		return float64(val), true
	case int32:
		val := value.(int32)
		return float64(val), true
	case int64:
		val := value.(int64)
		return float64(val), true
	case float32:
		val := value.(float32)
		return float64(val), true
	case float64:
		val := value.(float64)
		return float64(val), true
	default:
		fmt.Printf("The value is of an unknown type: %T\n"+"---"+key, t)
	}
	return 0, false
}

func (js *JSONObject) GetBool(key string) (bool, bool) {
	value, exists := js.Jm[key]
	if !exists {
		return false, false
	}
	switch t := value.(type) {
	case bool:
		val := value.(bool)
		return val, true
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return false, false
}
func (js *JSONObject) GetJSONObject(key string) JSONObject {
	val := js.Jm[key].(map[string]interface{})

	var jsonObj JSONObject
	jsonObj.Jm = val

	return jsonObj
}

func (js *JSONObject) GetJSONArray(key string) JSONArray {
	val := js.Jm[key].([]interface{})

	var jsonArray JSONArray
	jsonArray.Jm = val

	return jsonArray
}

func (jr *JSONArray) GetJSONObject(i int32) JSONObject {
	val := jr.Jm[i].(map[string]interface{})
	var jsonObj JSONObject
	jsonObj.Jm = val
	return jsonObj
}

func (jr *JSONArray) ToJsonString() string {
	data, _ := json.Marshal(jr.Jm)
	return string(data)
}

func (js *JSONObject) ToJsonString() string {
	data, _ := json.Marshal(js.Jm)
	return string(data)
}

func (jr *JSONArray) Size() int32 {
	if jr.Jm == nil {
		return 0
	}
	return int32(len(jr.Jm))
}

func (js *JSONObject) FluentPut(key string, value interface{}) *JSONObject {
	js.Jm[key] = value
	return js
}

func NewJSONObject() JSONObject {
	var jsonObj JSONObject
	jsonObj.Jm = make(map[string]interface{})
	return jsonObj
}

func (js *JSONObject) HasKey(key string) bool {
	_, exists := js.Jm[key]
	return exists
}
