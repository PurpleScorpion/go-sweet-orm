package mapper

import (
	"errors"

	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"gorm.io/gorm"
)

type UpdateWrapper struct {
	object       interface{}
	query        []queryCriteria
	updates      []updateSet
	groups       []*NestedQueryGroup // 添加嵌套查询组支持
	lastSQL      string
	txOrmer      *gorm.DB
	txFlag       bool
	excludeEmpty bool
	excludeField []string
	autoId       bool
	baseSql      string
	bulk         int
	values       []interface{}
}

func BuilderUpdateWrapper(flag bool) *UpdateWrapper {
	wrapper := &UpdateWrapper{}
	if flag {
		wrapper.txFlag = flag
		wrapper.txOrmer = globalDB.Begin()
	}
	wrapper.autoId = globalAutoId
	wrapper.bulk = 2000
	return wrapper
}

// 传播事务 , 可以将一个UpdateWrapper的事务传递到下一个UpdateWrapper中,并可以无限传播
// 若想手动传播可使用 SetTransaction 函数
func (qw *UpdateWrapper) SpreadTransaction(other UpdateWrapper) *UpdateWrapper {
	if !other.txFlag || other.txOrmer == nil {
		logger.Error("There are no transactions that can be disseminated")
	}
	qw.txFlag = true
	qw.txOrmer = other.txOrmer
	return qw
}

// 设置事务
func (qw *UpdateWrapper) SetTransaction(tx *gorm.DB) *UpdateWrapper {
	if tx == nil {
		logger.Error("The transaction parameter cannot be empty")
		return qw
	}
	qw.txFlag = true
	qw.txOrmer = tx
	return qw
}

func (qw *UpdateWrapper) SetBulk(bulk int) *UpdateWrapper {
	qw.bulk = bulk
	return qw
}

// 获取事务
func (qw *UpdateWrapper) GetTransaction() *gorm.DB {
	if qw.txFlag {
		return qw.txOrmer
	}
	return nil
}

// 手动开启事务
func (qw *UpdateWrapper) BeginTransaction() {
	if qw.txFlag && qw.txOrmer != nil {
		return
	}
	qw.txFlag = true
	qw.txOrmer = globalDB.Begin()
}

// 手动提交事务
func (qw *UpdateWrapper) Commit() error {
	if !qw.txFlag {
		return errors.New("Transaction not enabled")
	}
	if qw.txOrmer == nil {
		return errors.New("orm is nil")
	}
	qw.txOrmer.Commit()
	return nil
}

// 手动回滚事务
func (qw *UpdateWrapper) Rollback() error {
	if !qw.txFlag {
		return errors.New("Transaction not enabled")
	}
	if qw.txOrmer == nil {
		return errors.New("orm is nil")
	}
	qw.txOrmer.Rollback()
	return nil
}

func (qw *UpdateWrapper) CloseAutoId() *UpdateWrapper {
	qw.autoId = false
	return qw
}
func (qw *UpdateWrapper) OpenAutoId() *UpdateWrapper {
	qw.autoId = true
	return qw
}

func (qw *UpdateWrapper) SQL(sql string, values ...interface{}) *UpdateWrapper {
	qw.baseSql = sql
	qw.values = values
	return qw
}

func (qw *UpdateWrapper) SetExcludeEmpty(flag bool) *UpdateWrapper {
	qw.excludeEmpty = flag
	return qw
}

func (qw *UpdateWrapper) SetExcludeField(excludeField ...string) *UpdateWrapper {
	qw.excludeField = append(qw.excludeField, excludeField...)
	return qw
}

func (qw *UpdateWrapper) Set(flag bool, column string, value interface{}) *UpdateWrapper {
	checkParma(value)
	qw.updates = append(qw.updates, updateSet{condition: flag, columns: column, values: value})
	return qw
}

// 添加And方法以支持嵌套查询
func (qw *UpdateWrapper) And(group *QueryWrapper) *UpdateWrapper {
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

// 添加Or方法以支持嵌套查询
func (qw *UpdateWrapper) Or(group *QueryWrapper) *UpdateWrapper {
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

func (qw *UpdateWrapper) Eq(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, qw, "EQ")
}

func (qw *UpdateWrapper) Ne(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, qw, "NE")
}

func (qw *UpdateWrapper) Gt(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, qw, "GT")
}

func (qw *UpdateWrapper) Ge(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, qw, "GE")
}

func (qw *UpdateWrapper) Lt(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, qw, "LT")
}

func (qw *UpdateWrapper) Le(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, qw, "LE")
}
func (qw *UpdateWrapper) Like(flag bool, column string, value interface{}) *UpdateWrapper {
	val := likeValue(value, "LIKE")
	return addCondition4Update(flag, column, val, qw, "LIKE")
}

func (qw *UpdateWrapper) LikeLeft(flag bool, column string, value interface{}) *UpdateWrapper {
	val := likeValue(value, "LIKE_LEFT")
	return addCondition4Update(flag, column, val, qw, "LIKE")
}

func (qw *UpdateWrapper) LikeRight(flag bool, column string, value interface{}) *UpdateWrapper {
	val := likeValue(value, "LIKE_RIGHT")
	return addCondition4Update(flag, column, val, qw, "LIKE")
}

func (qw *UpdateWrapper) NotLike(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, qw, "NOT_LIKE")
}

func (qw *UpdateWrapper) IsNull(flag bool, column string) *UpdateWrapper {
	return addCondition4Update(flag, column, "", qw, "IS_NULL")
}

func (qw *UpdateWrapper) IsNotNull(flag bool, column string) *UpdateWrapper {
	return addCondition4Update(flag, column, "", qw, "IS_NOT_NULL")
}

func (qw *UpdateWrapper) In(flag bool, column string, values ...interface{}) *UpdateWrapper {
	checkParmas(values)
	return addConditionVals4Update(flag, column, values, qw, "IN")
}

func (qw *UpdateWrapper) InInt32(flag bool, column string, values []int32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "IN")
}

func (qw *UpdateWrapper) InInt64(flag bool, column string, values []int64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "IN")
}

func (qw *UpdateWrapper) InInt(flag bool, column string, values []int) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "IN")
}

func (qw *UpdateWrapper) InString(flag bool, column string, values []string) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "IN")
}

func (qw *UpdateWrapper) InFloat32(flag bool, column string, values []float32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "IN")
}

func (qw *UpdateWrapper) InFloat64(flag bool, column string, values []float64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "IN")
}

func (qw *UpdateWrapper) NotIn(flag bool, column string, values ...interface{}) *UpdateWrapper {
	checkParmas(values)
	return addConditionVals4Update(flag, column, values, qw, "NOT_IN")
}

func (qw *UpdateWrapper) NotInInt(flag bool, column string, values []int) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *UpdateWrapper) NotInInt32(flag bool, column string, values []int32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *UpdateWrapper) NotInInt64(flag bool, column string, values []int64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *UpdateWrapper) NotInString(flag bool, column string, values []string) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *UpdateWrapper) NotInFloat32(flag bool, column string, values []float32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *UpdateWrapper) NotInFloat64(flag bool, column string, values []float64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, qw, "NOT_IN")
}

func (qw *UpdateWrapper) Between(flag bool, column string, from string, to string) *UpdateWrapper {
	values := []interface{}{from, to}
	return addConditionVals4Update(flag, column, values, qw, "BETWEEN")
}

func (qw *UpdateWrapper) LastSql(sql string) *UpdateWrapper {
	qw.lastSQL = sql
	return qw
}

func addCondition4Update(flag bool, column string, value interface{}, q *UpdateWrapper, sqlKeyword string) *UpdateWrapper {
	checkParma(value)
	q.query = append(q.query, queryCriteria{
		condition: flag,
		columns:   column,
		values:    []interface{}{value},
		actions:   sqlKeyword,
	})
	return q
}

func addConditionVals4Update(flag bool, column string, value []interface{}, q *UpdateWrapper, sqlKeyword string) *UpdateWrapper {
	q.query = append(q.query, queryCriteria{
		condition: flag,
		columns:   column,
		values:    value,
		actions:   sqlKeyword,
	})
	return q
}
