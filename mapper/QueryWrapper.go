package mapper

import (
	"fmt"
	"reflect"
	"strconv"
)

type QueryWrapper struct {
	resList interface{}
	query   []queryCriteria
	sorts   []querySort
	updates []updateSet
	lastSQL string
}

type queryCriteria struct {
	condition bool
	columns   string
	values    []interface{}
	actions   string
}

type updateSet struct {
	condition bool
	columns   string
	values    interface{}
}

/*
condition 是否执行此查询

	true: 执行
	false: 不执行

isAsc

	true : asc
	false : desc
*/
type querySort struct {
	condition bool
	isAsc     bool
	columns   string
}

func BuilderQueryWrapper(list interface{}) QueryWrapper {
	var qw QueryWrapper
	qw.resList = list
	return qw
}

func (q *QueryWrapper) Set(flag bool, column string, value interface{}) *QueryWrapper {
	checkParma(value)
	q.updates = append(q.updates, updateSet{condition: flag, columns: column, values: value})
	return q
}

func (q *QueryWrapper) Eq(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, q, "EQ")
}

func (q *QueryWrapper) Ne(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, q, "NE")
}

func (q *QueryWrapper) Gt(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, q, "GT")
}

func (q *QueryWrapper) Ge(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, q, "GE")
}

func (q *QueryWrapper) Lt(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, q, "LT")
}

func (q *QueryWrapper) Le(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, q, "LE")
}
func (q *QueryWrapper) Like(flag bool, column string, value interface{}) *QueryWrapper {
	val := likeValue(value, "LIKE")
	return addCondition(flag, column, val, q, "LIKE")
}

func (q *QueryWrapper) LikeLeft(flag bool, column string, value interface{}) *QueryWrapper {
	val := likeValue(value, "LIKE_LEFT")
	return addCondition(flag, column, val, q, "LIKE")
}

func (q *QueryWrapper) LikeRight(flag bool, column string, value interface{}) *QueryWrapper {
	val := likeValue(value, "LIKE_RIGHT")
	return addCondition(flag, column, val, q, "LIKE")
}

func (q *QueryWrapper) NotLike(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, q, "NOT_LIKE")
}

func (q *QueryWrapper) IsNull(flag bool, column string) *QueryWrapper {
	return addCondition(flag, column, "", q, "IS_NULL")
}

func (q *QueryWrapper) IsNotNull(flag bool, column string) *QueryWrapper {
	return addCondition(flag, column, "", q, "IS_NOT_NULL")
}

func (q *QueryWrapper) In(flag bool, column string, values []interface{}) *QueryWrapper {
	checkParmas(values)
	return addConditionVals(flag, column, values, q, "IN")
}
func (q *QueryWrapper) NotIn(flag bool, column string, values []interface{}) *QueryWrapper {
	checkParmas(values)
	return addConditionVals(flag, column, values, q, "NOT_IN")
}

func (q *QueryWrapper) Between(flag bool, column string, from string, to string) *QueryWrapper {
	values := []interface{}{from, to}
	return addConditionVals(flag, column, values, q, "BETWEEN")
}

func (q *QueryWrapper) OrderByTimeAsc(flag bool, column string) *QueryWrapper {
	q.sorts = append(q.sorts, querySort{
		condition: flag,
		isAsc:     true,
		columns:   changeTimeData(column),
	})
	return q
}

func (q *QueryWrapper) OrderByAsc(flag bool, column string) *QueryWrapper {
	q.sorts = append(q.sorts, querySort{
		condition: flag,
		isAsc:     true,
		columns:   column,
	})
	return q
}

func (q *QueryWrapper) OrderByTimeDesc(flag bool, column string) *QueryWrapper {
	q.sorts = append(q.sorts, querySort{
		condition: flag,
		isAsc:     false,
		columns:   changeTimeData(column),
	})
	return q
}
func (q *QueryWrapper) OrderByDesc(flag bool, column string) *QueryWrapper {
	q.sorts = append(q.sorts, querySort{
		condition: flag,
		isAsc:     false,
		columns:   column,
	})
	return q
}

func (q *QueryWrapper) LastSql(sql string) *QueryWrapper {
	q.lastSQL = sql
	return q
}

func addCondition(flag bool, column string, value interface{}, q *QueryWrapper, sqlKeyword string) *QueryWrapper {
	checkParma(value)
	q.query = append(q.query, queryCriteria{
		condition: flag,
		columns:   column,
		values:    []interface{}{value},
		actions:   sqlKeyword,
	})
	return q
}
func addConditionVals(flag bool, column string, value []interface{}, q *QueryWrapper, sqlKeyword string) *QueryWrapper {
	q.query = append(q.query, queryCriteria{
		condition: flag,
		columns:   column,
		values:    value,
		actions:   sqlKeyword,
	})
	return q
}

func checkParmas(values []interface{}) {
	if values == nil {
		panic("The parameter type must be an integer, float, or string type")
	}
	if len(values) == 0 {
		panic("The parameter length cannot be 0")
	}
	one := values[0]
	checkParma(one)
	if len(values) == 1 {
		return
	}
	oneType := reflect.TypeOf(one)

	for i := 1; i < len(values); i++ {
		if reflect.TypeOf(values[i]) != oneType {
			panic("The parameter type must be the same")
		}
	}

}

func likeValue(value interface{}, sqlKeyword string) interface{} {
	if value == nil {
		panic("The parameter type must be an integer, float, or string type")
	}
	var valStr string
	switch value.(type) {
	case int:
		valStr = int64ToString(int64(value.(int)))
	case int8:
		valStr = int64ToString(int64(value.(int8)))
	case int16:
		valStr = int64ToString(int64(value.(int16)))
	case int32:
		valStr = int64ToString(int64(value.(int32)))
	case int64:
		valStr = int64ToString(value.(int64))
	case uint:
		valStr = int64ToString(int64(value.(uint)))
	case uint8:
		valStr = int64ToString(int64(value.(uint8)))
	case uint16:
		valStr = int64ToString(int64(value.(uint16)))
	case uint32:
		valStr = int64ToString(int64(value.(uint32)))
	case uint64:
		valStr = int64ToString(int64(value.(uint64)))
	case float32:
		valStr = float64ToString(float64(value.(float32)))
	case float64:
		valStr = float64ToString(value.(float64))
	case string:
		valStr = value.(string)
	default:
		// 都不是，返回错误
		panic("The parameter type must be an integer, float, or string type")
	}

	if sqlKeyword == "LIKE_LEFT" {
		return "%" + valStr
	} else if sqlKeyword == "LIKE_RIGHT" {
		return valStr + "%"
	} else if sqlKeyword == "LIKE" {
		return "%" + valStr + "%"
	}
	panic("Unknown SqlKeyword (" + sqlKeyword + ") appears")
}

func checkParma(value interface{}) {
	if value == nil {
		panic("Parameter cannot be empty")
	}
	switch tp := value.(type) {
	case int:
		return
	case int8:
		return
	case int16:
		return
	case int32:
		return
	case int64:
		return
	case uint:
		return
	case uint8:
		return
	case uint16:
		return
	case uint32:
		return
	case uint64:
		return
	case float32:
		return
	case float64:
		return
	case bool:
		return
	case string:
		return
	default:
		// 都不是，返回错误
		panic(fmt.Sprintf("The parameter type must be an integer, float, or string type, but the type is %v", tp))
	}
}

func int64ToString(num int64) string {
	str := strconv.FormatInt(num, 10)
	return str
}

func float64ToString(floatValue float64) string {
	str := strconv.FormatFloat(floatValue, 'f', -1, 64)
	return str
}

func changeTimeData(val string) string {
	if ActiveDB == Sqlite {
		val = fmt.Sprintf("datetime(%s)", val)
		return val
	}
	return val
}
