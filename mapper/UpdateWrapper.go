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

func BuilderUpdateWrapper(obj interface{}, flag ...bool) UpdateWrapper {
	var wrapper UpdateWrapper
	if len(flag) == 0 {
		if !transaction {
			wrapper.txFlag = false
		} else {
			wrapper.txFlag = true
		}
	} else {
		wrapper.txFlag = flag[0]
	}

	if wrapper.txFlag {
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
func (qw *UpdateWrapper) SetTransaction(tx orm.TxOrmer) *UpdateWrapper {
	if tx == nil {
		logger.Error("The transaction parameter cannot be empty")
		return qw
	}
	qw.txFlag = true
	qw.txOrmer = tx
	return qw
}

// 获取事务
func (qw *UpdateWrapper) GetTransaction() orm.TxOrmer {
	if qw.txFlag {
		return qw.txOrmer
	}
	return nil
}

// 手动开启事务
func (qw *UpdateWrapper) BenginTransaction() error {
	if qw.txFlag {
		if qw.txOrmer != nil {
			return nil
		}
	}
	qw.txFlag = true
	txOrmer, err := GetOrm().Begin()
	if err != nil {
		qw.txFlag = false
		logger.Error("Transaction initiation failed: %v", err)
		return err
	}
	qw.txOrmer = txOrmer
	return nil
}

// 手动提交事务
func (qw *UpdateWrapper) Commit() error {
	if !qw.txFlag {
		return errors.New("Transaction not enabled")
	}
	if qw.txOrmer == nil {
		return errors.New("orm is nil")
	}
	return qw.txOrmer.Commit()
}

// 手动回滚事务
func (qw *UpdateWrapper) Rollback() error {
	if !qw.txFlag {
		return errors.New("Transaction not enabled")
	}
	if qw.txOrmer == nil {
		return errors.New("orm is nil")
	}
	return qw.txOrmer.Rollback()
}

func (qw *UpdateWrapper) Set(flag bool, column string, value interface{}) *UpdateWrapper {
	checkParma(value)
	qw.updates = append(qw.updates, updateSet{condition: flag, columns: column, values: value})
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
