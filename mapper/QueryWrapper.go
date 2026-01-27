package mapper

import (
	"fmt"
	"reflect"
	"strconv"
)

type QueryWrapper struct {
	resList interface{}
	query   []queryCriteria
	groups  []*NestedQueryGroup // 支持嵌套查询组
	sorts   []querySort
	lastSQL string
	groupType string // 用于标记这是一个AND组还是OR组，用于内部连接符
}

type queryCriteria struct {
	condition bool
	columns   string
	values    []interface{}
	actions   string
}

// NestedQueryGroup represents a group of nested query conditions
type NestedQueryGroup struct {
	groupType string              // "AND" or "OR"
	criteria  []queryCriteria     // individual criteria in this group
	groups    []*NestedQueryGroup // nested groups inside this group
	actions   string              // "GROUP_AND" or "GROUP_OR"
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

func BuilderQueryWrapper() *QueryWrapper {
	wrapper := &QueryWrapper{}
	return wrapper
}

// 为内部创建使用的新构造函数
func newQueryGroup(groupType string) *QueryWrapper {
	return &QueryWrapper{
		groupType: groupType,
	}
}

// NewAndGroup creates a new AND group for nested queries
func NewAndGroup() *QueryWrapper {
	return newQueryGroup("AND")
}

// NewOrGroup creates a new OR group for nested queries
func NewOrGroup() *QueryWrapper {
	return newQueryGroup("OR")
}

func (qw *QueryWrapper) And(group *QueryWrapper) *QueryWrapper {
	if group == nil || (len(group.query) == 0 && len(group.groups) == 0) {
		return qw
	}

	// Create a nested group - use the group's own type if defined, otherwise default to AND
	groupType := "AND"
	if group.groupType != "" {
		groupType = group.groupType
	}
	
	groupContainer := &NestedQueryGroup{
		groupType: groupType,
		criteria:  group.query,
		groups:    group.groups,
		actions:   "GROUP_AND",
	}
	qw.groups = append(qw.groups, groupContainer)
	return qw
}

func (qw *QueryWrapper) Or(group *QueryWrapper) *QueryWrapper {
	if group == nil || (len(group.query) == 0 && len(group.groups) == 0) {
		return qw
	}

	// Create a nested group - use the group's own type if defined, otherwise default to OR
	groupType := "OR"
	if group.groupType != "" {
		groupType = group.groupType
	}
	
	groupContainer := &NestedQueryGroup{
		groupType: groupType,
		criteria:  group.query,
		groups:    group.groups,
		actions:   "GROUP_OR",
	}
	qw.groups = append(qw.groups, groupContainer)
	return qw
}

func (qw *QueryWrapper) Eq(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, qw, "EQ")
}

func (qw *QueryWrapper) Ne(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, qw, "NE")
}

func (qw *QueryWrapper) Gt(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, qw, "GT")
}

func (qw *QueryWrapper) Ge(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, qw, "GE")
}

func (qw *QueryWrapper) Lt(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, qw, "LT")
}

func (qw *QueryWrapper) Le(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, qw, "LE")
}
func (qw *QueryWrapper) Like(flag bool, column string, value interface{}) *QueryWrapper {
	val := likeValue(value, "LIKE")
	return addCondition(flag, column, val, qw, "LIKE")
}

func (qw *QueryWrapper) LikeLeft(flag bool, column string, value interface{}) *QueryWrapper {
	val := likeValue(value, "LIKE_LEFT")
	return addCondition(flag, column, val, qw, "LIKE")
}

func (qw *QueryWrapper) LikeRight(flag bool, column string, value interface{}) *QueryWrapper {
	val := likeValue(value, "LIKE_RIGHT")
	return addCondition(flag, column, val, qw, "LIKE")
}

func (qw *QueryWrapper) NotLike(flag bool, column string, value interface{}) *QueryWrapper {
	return addCondition(flag, column, value, qw, "NOT_LIKE")
}

func (qw *QueryWrapper) IsNull(flag bool, column string) *QueryWrapper {
	return addCondition(flag, column, "", qw, "IS_NULL")
}

func (qw *QueryWrapper) IsNotNull(flag bool, column string) *QueryWrapper {
	return addCondition(flag, column, "", qw, "IS_NOT_NULL")
}

func (qw *QueryWrapper) In(flag bool, column string, values ...interface{}) *QueryWrapper {
	checkParmas(values)
	return addConditionVals(flag, column, values, qw, "IN")
}

func (qw *QueryWrapper) InInt32(flag bool, column string, values []int32) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "IN")
}

func (qw *QueryWrapper) InInt64(flag bool, column string, values []int64) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "IN")
}

func (qw *QueryWrapper) InInt(flag bool, column string, values []int) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "IN")
}

func (qw *QueryWrapper) InString(flag bool, column string, values []string) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "IN")
}

func (qw *QueryWrapper) InFloat32(flag bool, column string, values []float32) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "IN")
}

func (qw *QueryWrapper) InFloat64(flag bool, column string, values []float64) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "IN")
}

func (qw *QueryWrapper) NotIn(flag bool, column string, values ...interface{}) *QueryWrapper {
	checkParmas(values)
	return addConditionVals(flag, column, values, qw, "NOT_IN")
}

func (qw *QueryWrapper) NotInInt(flag bool, column string, values []int) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *QueryWrapper) NotInInt32(flag bool, column string, values []int32) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *QueryWrapper) NotInInt64(flag bool, column string, values []int64) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *QueryWrapper) NotInString(flag bool, column string, values []string) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *QueryWrapper) NotInFloat32(flag bool, column string, values []float32) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *QueryWrapper) NotInFloat64(flag bool, column string, values []float64) *QueryWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *QueryWrapper) Between(flag bool, column string, from string, to string) *QueryWrapper {
	values := []interface{}{from, to}
	return addConditionVals(flag, column, values, qw, "BETWEEN")
}

func (qw *QueryWrapper) OrderByTimeAsc(flag bool, column string) *QueryWrapper {
	qw.sorts = append(qw.sorts, querySort{
		condition: flag,
		isAsc:     true,
		columns:   changeTimeData(column),
	})
	return qw
}

func (qw *QueryWrapper) OrderByAsc(flag bool, column string) *QueryWrapper {
	qw.sorts = append(qw.sorts, querySort{
		condition: flag,
		isAsc:     true,
		columns:   column,
	})
	return qw
}

func (qw *QueryWrapper) OrderByTimeDesc(flag bool, column string) *QueryWrapper {
	qw.sorts = append(qw.sorts, querySort{
		condition: flag,
		isAsc:     false,
		columns:   changeTimeData(column),
	})
	return qw
}
func (qw *QueryWrapper) OrderByDesc(flag bool, column string) *QueryWrapper {
	qw.sorts = append(qw.sorts, querySort{
		condition: flag,
		isAsc:     false,
		columns:   column,
	})
	return qw
}

func (qw *QueryWrapper) LastSql(sql string) *QueryWrapper {
	qw.lastSQL = sql
	return qw
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
	if active_db == Sqlite {
		val = fmt.Sprintf("datetime(%s)", val)
		return val
	}
	return val
}
