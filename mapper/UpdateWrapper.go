package mapper

import (
	"errors"
	"github.com/PurpleScorpion/go-sweet-orm/logger"
	"github.com/beego/beego/v2/client/orm"
)

type UpdateWrapper struct {
	object  interface{}
	query   []queryCriteria
	updates []updateSet
	lastSQL string
	txOrmer orm.TxOrmer
	txFlag  bool
}

func BuilderUpdateWrapper(obj interface{}, flag bool) UpdateWrapper {
	var wrapper UpdateWrapper
	wrapper.txFlag = flag
	if flag {
		txOrmer, err := GetOrm().Begin()
		if err != nil {
			logger.Error("Transaction initiation failed: %v", err)
			wrapper.txFlag = false
		} else {
			wrapper.txOrmer = txOrmer
		}
	}
	wrapper.object = obj
	return wrapper
}

func (q *UpdateWrapper) GetTransaction() orm.TxOrmer {
	if q.txFlag {
		return q.txOrmer
	}
	return nil
}

func (q *UpdateWrapper) BenginTransaction() error {
	if q.txFlag {
		if q.txOrmer != nil {
			return nil
		}
	}
	q.txFlag = true
	txOrmer, err := GetOrm().Begin()
	if err != nil {
		q.txFlag = false
		logger.Error("Transaction initiation failed: %v", err)
		return err
	}
	q.txOrmer = txOrmer
	return nil
}

func (q *UpdateWrapper) Commit() error {
	if !q.txFlag {
		return errors.New("Transaction not enabled")
	}
	if q.txOrmer == nil {
		return errors.New("orm is nil")
	}
	return q.txOrmer.Commit()
}

func (q *UpdateWrapper) Rollback() error {
	if !q.txFlag {
		return errors.New("Transaction not enabled")
	}
	if q.txOrmer == nil {
		return errors.New("orm is nil")
	}
	return q.txOrmer.Rollback()
}

func (q *UpdateWrapper) Set(flag bool, column string, value interface{}) *UpdateWrapper {
	checkParma(value)
	q.updates = append(q.updates, updateSet{condition: flag, columns: column, values: value})
	return q
}

func (q *UpdateWrapper) Eq(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, q, "EQ")
}

func (q *UpdateWrapper) Ne(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, q, "NE")
}

func (q *UpdateWrapper) Gt(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, q, "GT")
}

func (q *UpdateWrapper) Ge(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, q, "GE")
}

func (q *UpdateWrapper) Lt(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, q, "LT")
}

func (q *UpdateWrapper) Le(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, q, "LE")
}
func (q *UpdateWrapper) Like(flag bool, column string, value interface{}) *UpdateWrapper {
	val := likeValue(value, "LIKE")
	return addCondition4Update(flag, column, val, q, "LIKE")
}

func (q *UpdateWrapper) LikeLeft(flag bool, column string, value interface{}) *UpdateWrapper {
	val := likeValue(value, "LIKE_LEFT")
	return addCondition4Update(flag, column, val, q, "LIKE")
}

func (q *UpdateWrapper) LikeRight(flag bool, column string, value interface{}) *UpdateWrapper {
	val := likeValue(value, "LIKE_RIGHT")
	return addCondition4Update(flag, column, val, q, "LIKE")
}

func (q *UpdateWrapper) NotLike(flag bool, column string, value interface{}) *UpdateWrapper {
	return addCondition4Update(flag, column, value, q, "NOT_LIKE")
}

func (q *UpdateWrapper) IsNull(flag bool, column string) *UpdateWrapper {
	return addCondition4Update(flag, column, "", q, "IS_NULL")
}

func (q *UpdateWrapper) IsNotNull(flag bool, column string) *UpdateWrapper {
	return addCondition4Update(flag, column, "", q, "IS_NOT_NULL")
}

func (q *UpdateWrapper) In(flag bool, column string, values ...interface{}) *UpdateWrapper {
	checkParmas(values)
	return addConditionVals4Update(flag, column, values, q, "IN")
}

func (q *UpdateWrapper) InInt32(flag bool, column string, values []int32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "IN")
}

func (q *UpdateWrapper) InInt64(flag bool, column string, values []int64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "IN")
}

func (q *UpdateWrapper) InInt(flag bool, column string, values []int) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "IN")
}

func (q *UpdateWrapper) InString(flag bool, column string, values []string) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "IN")
}

func (q *UpdateWrapper) InFloat32(flag bool, column string, values []float32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "IN")
}

func (q *UpdateWrapper) InFloat64(flag bool, column string, values []float64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "IN")
}

func (q *UpdateWrapper) NotIn(flag bool, column string, values ...interface{}) *UpdateWrapper {
	checkParmas(values)
	return addConditionVals4Update(flag, column, values, q, "NOT_IN")
}

func (q *UpdateWrapper) NotInInt(flag bool, column string, values []int) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "NOT_IN")
}

func (q *UpdateWrapper) NotInInt32(flag bool, column string, values []int32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "NOT_IN")
}

func (q *UpdateWrapper) NotInInt64(flag bool, column string, values []int64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "NOT_IN")
}

func (q *UpdateWrapper) NotInString(flag bool, column string, values []string) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "NOT_IN")
}

func (q *UpdateWrapper) NotInFloat32(flag bool, column string, values []float32) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "NOT_IN")
}

func (q *UpdateWrapper) NotInFloat64(flag bool, column string, values []float64) *UpdateWrapper {
	var interfaces []interface{}
	for _, s := range values {
		interfaces = append(interfaces, s)
	}
	return addConditionVals4Update(flag, column, interfaces, q, "NOT_IN")
}

func (q *UpdateWrapper) Between(flag bool, column string, from string, to string) *UpdateWrapper {
	values := []interface{}{from, to}
	return addConditionVals4Update(flag, column, values, q, "BETWEEN")
}

func (q *UpdateWrapper) LastSql(sql string) *UpdateWrapper {
	q.lastSQL = sql
	return q
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
