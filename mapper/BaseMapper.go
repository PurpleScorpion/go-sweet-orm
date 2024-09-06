package mapper

import (
	"errors"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-orm/v2/logger"
	"github.com/beego/beego/v2/client/orm"
	"reflect"
	"strings"
)

type BaseMapper struct {
}

func Page(page PageUtils) PageData {
	qw := page.wrapper
	offSet := page.getOffSet()
	pageSize := page.getPageSize()
	qw.lastSQL = fmt.Sprintf("limit %d,%d", offSet, pageSize)
	err := qw.SelectList()
	if err != nil {
		panic(err)
	}
	qw.lastSQL = ""
	count := qw.SelectCount()
	page.setTotalSize(count)
	pageData := page.pageData()
	return pageData
}

/**
 * 根据主键查询
 * 该查询会根据传入的参数进行两种结果类型封装
 * 1. 当传入为对象类型时, 会将结果集封装为一个对象(不推荐,因为go的空值没有null这个结果,只能判断id是否为0)
 * 2. 当传入对象为数组类型时 , 会将结果集封装为一个数组(推荐,可以通过判断数组length是否为0来判断是否存在数据)
 */
func SelectById(obj interface{}, id interface{}) error {
	idFieldName := getTableId(obj)
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}

	tableName := getTableName(obj)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	tp := getParmaStruct(obj)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}

	if o == nil {
		return errors.New("orm is null")
	}

	sql := fmt.Sprintf("select * from %s where %s = ?", tableName, idFieldName)
	LogInfo("SelectById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("SelectById", fmt.Sprintf("Parameters: %v", id))
	r := o.Raw(sql, id)
	if r == nil {
		return errors.New("RawSeter is null")
	}
	if tp == "Array" {
		_, err := r.QueryRows(obj)
		if err != nil {
			return err
		}
	} else {
		err := r.QueryRow(obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qw *QueryWrapper) SelectList() error {
	resList := qw.resList
	tp := getParmaStruct(resList)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	if tp != "Array" {
		panic("The result set parameter must be of array type")
	}
	tableName := getTableName(resList)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	baseSQL, values := qw.queryWrapper4SQL()

	baseSQL = fmt.Sprintf("select * from %s where 1=1 %s", tableName, baseSQL)
	LogInfo("SelectList", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("SelectList", fmt.Sprintf("Parameters: %v", values))
	if o == nil {
		return errors.New("orm is nil")
	}
	r := o.Raw(baseSQL, values...)
	if r == nil {
		return errors.New("RawSeter is null")
	}

	_, err := r.QueryRows(resList)
	if err != nil {
		return err
	}
	qw.resList = resList
	return nil
}

func (qw *QueryWrapper) SelectCount() int64 {
	resList := qw.resList
	tp := getParmaStruct(resList)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	tableName := getTableName(resList)
	if tableName == "null" {
		panic("TableName method does not exist")
	}
	baseSQL, values := getQuerySQL(qw.query)
	if !isEmpty(qw.lastSQL) {
		baseSQL = fmt.Sprintf("%s %s", baseSQL, qw.lastSQL)
	}
	baseSQL = fmt.Sprintf("select count(*) from %s where 1=1 %s", tableName, baseSQL)
	LogInfo("SelectCount", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("SelectCount", fmt.Sprintf("Parameters: %v", values))
	var count int64
	if o == nil {
		return 0
	}

	r := o.Raw(baseSQL, values...)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	err := r.QueryRow(&count)
	if err != nil {
		return 0
	}
	return count
}

func (qw *UpdateWrapper) Delete() int64 {
	resList := qw.object
	tp := getParmaStruct(resList)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	tableName := getTableName(resList)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	if o == nil {
		return 0
	}

	baseSQL := fmt.Sprintf("delete from %s where 1=1 ", tableName)
	// 查询组合器
	sql, values := getQuerySQL(qw.query)
	// 拼接sql
	baseSQL = fmt.Sprintf("%s %s ", baseSQL, sql)
	// 拼接lastsql
	if !isEmpty(qw.lastSQL) {
		baseSQL = fmt.Sprintf("%s %s", baseSQL, qw.lastSQL)
	}
	LogInfo("Delete", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("Delete", fmt.Sprintf("Parameters: %v", values))
	r := o.Raw(baseSQL, values...)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	res, err := r.Exec()
	if err != nil {
		logger.Error("Delete failed: %v", err)
		return 0
	}
	if res == nil {
		logger.Error("Delete exec result is null")
		return 0
	}

	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Delete RowsAffected failed: %v", err)
		return 0
	}
	return count
}

func (qw *UpdateWrapper) DeleteByIds(ids interface{}) int64 {
	return delByIds(qw.object, ids, qw)
}
func DeleteByIds(obj interface{}, ids interface{}) int64 {
	return delByIds(obj, ids, nil)
}

func delByIds(obj interface{}, ids interface{}, qw *UpdateWrapper) int64 {
	idFieldName := getTableId(obj)
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}

	tableName := getTableName(obj)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	tp := getParmaStruct(obj)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}

	if !isSlice(ids) {
		panic("ids parameter is not a slice")
	}
	if qw != nil {
		return delByIds4Tx(idFieldName, tableName, ids, qw)
	}
	return delByIds4NoTx(idFieldName, tableName, ids)
}

func delByIds4Tx(idFieldName, tableName string, ids interface{}, qw *UpdateWrapper) int64 {
	if !qw.txFlag || qw.txOrmer == nil {
		return delById4NoTx(idFieldName, tableName, ids)
	}

	// 获取 ids 的反射值
	idsVal := reflect.ValueOf(ids)
	var str strings.Builder
	str.WriteString("(")
	// 遍历切片
	for i := 0; i < idsVal.Len(); i++ {
		if i == idsVal.Len()-1 {
			str.WriteString("?")
		} else {
			str.WriteString("?,")
		}
	}
	str.WriteString(")")
	sql := fmt.Sprintf("delete from %s where %s in %s", tableName, idFieldName, str.String())
	LogInfo("DeleteByIds", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("DeleteByIds", fmt.Sprintf("Parameters: %v", ids))

	r := qw.txOrmer.Raw(sql, ids)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	res, err := r.Exec()
	if err != nil {
		logger.Error("Delete failed: %v", err)
		return 0
	}
	if res == nil {
		logger.Error("Delete exec result is null")
		return 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Delete RowsAffected failed: %v", err)
		return 0
	}
	return count
}

func delByIds4NoTx(idFieldName, tableName string, ids interface{}) int64 {
	if o == nil {
		return 0
	}
	// 获取 ids 的反射值
	idsVal := reflect.ValueOf(ids)
	var str strings.Builder
	str.WriteString("(")
	// 遍历切片
	for i := 0; i < idsVal.Len(); i++ {
		if i == idsVal.Len()-1 {
			str.WriteString("?")
		} else {
			str.WriteString("?,")
		}
	}
	str.WriteString(")")
	sql := fmt.Sprintf("delete from %s where %s in %s", tableName, idFieldName, str.String())
	LogInfo("DeleteById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("DeleteById", fmt.Sprintf("Parameters: %v", ids))

	r := o.Raw(sql, ids)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	res, err := r.Exec()

	if err != nil {
		logger.Error("Delete failed: %v", err)
		return 0
	}
	if res == nil {
		logger.Error("Delete exec result is null")
		return 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Delete RowsAffected failed: %v", err)
		return 0
	}
	return count
}

// 使用事务的删除
func (qw *UpdateWrapper) DeleteById(id interface{}) int64 {
	return delById(qw.object, id, qw)
}

// 不使用事务的快捷删除
func DeleteById(obj interface{}, id interface{}) int64 {
	return delById(obj, id, nil)
}

func delById(obj interface{}, id interface{}, qw *UpdateWrapper) int64 {
	idFieldName := getTableId(obj)
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}

	tableName := getTableName(obj)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	tp := getParmaStruct(obj)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}

	if qw != nil {
		return delById4Tx(idFieldName, tableName, id, qw)
	}
	return delById4NoTx(idFieldName, tableName, id)
}

func delById4Tx(idFieldName, tableName string, id interface{}, qw *UpdateWrapper) int64 {

	if !qw.txFlag || qw.txOrmer == nil {
		return delById4NoTx(idFieldName, tableName, id)
	}
	sql := fmt.Sprintf("delete from %s where %s = ?", tableName, idFieldName)
	LogInfo("DeleteById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("DeleteById", fmt.Sprintf("Parameters: %v", id))

	r := qw.txOrmer.Raw(sql, id)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	res, err := r.Exec()
	if err != nil {
		logger.Error("Delete failed: %v", err)
		return 0
	}
	if res == nil {
		logger.Error("Delete exec result is null")
		return 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Delete RowsAffected failed: %v", err)
		return 0
	}
	return count

}

func delById4NoTx(idFieldName, tableName string, id interface{}) int64 {
	if o == nil {
		return 0
	}
	sql := fmt.Sprintf("delete from %s where %s = ?", tableName, idFieldName)
	LogInfo("DeleteById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("DeleteById", fmt.Sprintf("Parameters: %v", id))

	r := o.Raw(sql, id)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	res, err := r.Exec()

	if err != nil {
		logger.Error("Delete failed: %v", err)
		return 0
	}
	if res == nil {
		logger.Error("Delete exec result is null")
		return 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Delete RowsAffected failed: %v", err)
		return 0
	}
	return count
}

func (qw *UpdateWrapper) Update() int64 {
	resList := qw.object
	tp := getParmaStruct(resList)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	tableName := getTableName(resList)
	if tableName == "null" {
		panic("TableName method does not exist")
	}
	if qw.updates == nil || len(qw.updates) == 0 {
		panic("Missing fields to be updated")
	}
	values := make([]interface{}, 0)
	baseSQL := fmt.Sprintf("update %s ", tableName)

	// 去除flage为false的更新字段
	qw.removeFalseUpdates()

	if len(qw.updates) == 1 {
		baseSQL = fmt.Sprintf("%s set %s = ? ", baseSQL, qw.updates[0].columns)
		values = append(values, qw.updates[0].values)
	} else {
		for i := 0; i < len(qw.updates); i++ {
			if !qw.updates[i].condition {
				continue
			}
			if i == len(qw.updates)-1 {
				baseSQL = fmt.Sprintf("%s %s = ? ", baseSQL, qw.updates[i].columns)
			} else if i == 0 {
				baseSQL = fmt.Sprintf("%s set %s = ?, ", baseSQL, qw.updates[i].columns)
			} else {
				baseSQL = fmt.Sprintf("%s %s = ?, ", baseSQL, qw.updates[i].columns)
			}
			values = append(values, qw.updates[i].values)
		}
	}

	baseSQL = fmt.Sprintf("%s where 1=1 ", baseSQL)
	sql, vals := getQuerySQL(qw.query)

	baseSQL = fmt.Sprintf("%s %s ", baseSQL, sql)
	if !isEmpty(qw.lastSQL) {
		baseSQL = fmt.Sprintf("%s %s ", baseSQL, qw.lastSQL)
	}

	for i := 0; i < len(vals); i++ {
		values = append(values, vals[i])
	}
	LogInfo("Update", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("Update", fmt.Sprintf("Parameters: %v", values))

	var r orm.RawSeter

	if !qw.txFlag || qw.txOrmer == nil {
		if o == nil {
			return 0
		}
		r = o.Raw(baseSQL, values...)
	} else {
		r = qw.txOrmer.Raw(baseSQL, values...)
	}
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	res, err := r.Exec()
	if err != nil {
		logger.Error("Update failed: %v", err)
		return 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Update RowsAffected failed: %v", err)
		return 0
	}
	return count

}

// 批量插入 - 可极大提高效率
// bulk, 单次插入数量
func (qw *UpdateWrapper) InsertAll(bulk int, excludeEmpty bool, excludeField ...string) int64 {
	sind := reflect.Indirect(reflect.ValueOf(qw.object))
	switch sind.Kind() {
	case reflect.Array, reflect.Slice:
		if sind.Len() == 0 {
			return 0
		}
	default:
		return 0
	}

	if bulk <= 1 {
		var count int64 = 0
		for i := 0; i < sind.Len(); i++ {
			count += baseInsert(sind.Index(i).Interface(), true, qw, excludeField...)
		}
		return count
	} else {
		idFieldName := getTableId(sind.Index(0).Interface())
		if idFieldName == "null" {
			panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
		}

		tableName := getTableName(sind.Index(0).Interface())
		if tableName == "null" {
			panic("TableName method does not exist")
		}

		return insertMulti(idFieldName, tableName, qw, bulk, excludeEmpty, excludeField...)
	}
}

func insertMulti(idFieldName, tableName string, qw *UpdateWrapper, bulk int, excludeEmpty bool, excludeField ...string) int64 {
	sind := reflect.ValueOf(qw.object)
	fields, _, _ := getExcludeFiledName(sind.Index(0).Interface(), tableName, excludeField, idFieldName, excludeEmpty)
	if len(fields) == 0 {
		panic("No fields were found")
	}

	sql1 := "("
	for i := 0; i < len(fields); i++ {
		if i == len(fields)-1 {
			sql1 = fmt.Sprintf("%s %s", sql1, fields[i])
		} else {
			sql1 = fmt.Sprintf("%s %s, ", sql1, fields[i])
		}
	}
	sql1 = fmt.Sprintf("%s ) ", sql1)
	endOffset := 1
	if bulk > sind.Len() {
		bulk = sind.Len()
	} else {
		endOffset = (sind.Len() / bulk) + 1
	}
	baseSQL := fmt.Sprintf("insert into %s %s values ", tableName, sql1)
	var sum int64 = 0
	// 遍历偏移量
	for offset := 0; offset < endOffset; offset++ {
		startIndex := offset * bulk
		endIndex := offset*bulk + bulk
		if endIndex > sind.Len() {
			endIndex = sind.Len()
		}
		// 遍历所有数据
		ss := ""
		values := make([]interface{}, 0)
		for i := startIndex; i < endIndex; i++ {
			sql2 := "("
			fs, vals, _ := getExcludeFiledName(sind.Index(i).Interface(), tableName, excludeField, idFieldName, excludeEmpty)
			//遍历字段
			for x := 0; x < len(fs); x++ {

				if x == len(fs)-1 {
					sql2 = fmt.Sprintf("%s ?", sql2)
				} else {
					sql2 = fmt.Sprintf("%s ?, ", sql2)
				}
			}
			if i == endIndex-1 {
				sql2 = fmt.Sprintf("%s ) ", sql2)
			} else {
				sql2 = fmt.Sprintf("%s ), ", sql2)
			}
			ss += sql2
			values = append(values, vals...)
		}

		var r orm.RawSeter

		if !qw.txFlag || qw.txOrmer == nil {
			if o == nil {
				return 0
			}
			r = o.Raw(baseSQL+ss, values...)
		} else {
			r = qw.txOrmer.Raw(baseSQL+ss, values...)
		}
		if r == nil {
			logger.Error("RawSeter is null")
			return 0
		}
		res, err := r.Exec()
		if err != nil {
			logger.Error("Insert failed: %v", err)
			return 0
		}

		count, err := res.RowsAffected()
		if err != nil {
			logger.Error("Insert RowsAffected failed: %v", err)
			return 0
		}
		sum += count

	}
	return sum
}

func (qw *UpdateWrapper) Insert(excludeField ...string) int64 {
	return baseInsert(qw.object, true, qw, excludeField...)
}

// 默认 自增主键和空值排除的新增
func Insert(pojo interface{}, excludeField ...string) int64 {
	return baseInsert(pojo, true, nil, excludeField...)
}

func (qw *UpdateWrapper) InsertCustom(excludeEmpty bool, excludeField ...string) int64 {
	return baseInsert(qw.object, excludeEmpty, qw, excludeField...)
}

/*
pojo: 实体对象
autoId: 是否为自增主键
excludeEmpty: 是否排除空值
excludeField: 排除字段 - 数据库表字段名
*/
func baseInsert(pojo interface{}, excludeEmpty bool, qw *UpdateWrapper, excludeField ...string) int64 {
	idFieldName := getTableId(pojo)
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}

	tableName := getTableName(pojo)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	tp := getParmaStruct(pojo)
	if tp == "null" {
		panic("Input is neither a struct")
	}
	if tp == "slice" {
		panic("The pojo parameter cannot be of array type")
	}

	fields, values, autoId := getExcludeFiledName(pojo, tableName, excludeField, idFieldName, excludeEmpty)
	if len(fields) == 0 {
		panic("No fields were found")
	}

	str1 := "("
	str2 := "("
	for i := 0; i < len(fields); i++ {
		if i == len(fields)-1 {
			str1 = fmt.Sprintf("%s %s", str1, fields[i])
			str2 = fmt.Sprintf("%s ?", str2)
		} else {
			str1 = fmt.Sprintf("%s %s, ", str1, fields[i])
			str2 = fmt.Sprintf("%s ?, ", str2)
		}
	}
	str1 = fmt.Sprintf("%s ) ", str1)
	str2 = fmt.Sprintf("%s ) ", str2)

	baseSQL := fmt.Sprintf("insert into %s %s values %s", tableName, str1, str2)
	LogInfo("Insert", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("Insert", fmt.Sprintf("Parameters: %v", values))
	var r orm.RawSeter

	if qw == nil || !qw.txFlag || qw.txOrmer == nil {
		if o == nil {
			return 0
		}
		r = o.Raw(baseSQL, values...)
	} else {
		r = qw.txOrmer.Raw(baseSQL, values...)
	}
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	res, err := r.Exec()
	if err != nil {
		logger.Error("Insert failed: %v", err)
		return 0
	}
	if autoId {
		lastId, err1 := res.LastInsertId()
		if err1 != nil {
			logger.Error("Get LastInsertId failed: %v", err1)
			return 0
		}
		saveLastInsertId(pojo, lastId)
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Insert RowsAffected failed: %v", err)
		return 0
	}
	return count
}

func SelectCount4SQL(sql string, values ...interface{}) int64 {
	var count int64
	if o == nil {
		return 0
	}
	r := o.Raw(sql, values...)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}
	err := r.QueryRow(&count)
	if err != nil {
		return 0
	}
	return count
}
func SelectList4SQL(resList interface{}, sql string, values ...interface{}) error {
	tp := getParmaStruct(resList)
	if tp == "null" {
		return errors.New("Input is neither a slice")
	}
	if tp != "Array" {
		return errors.New("The result set parameter must be of array type")
	}
	if o == nil {
		return errors.New("orm is nil")
	}
	r := o.Raw(sql, values...)
	if r == nil {
		return errors.New("RawSeter is null")
	}
	_, err := r.QueryRows(resList)
	if err != nil {
		return err
	}
	return nil
}

/*
原生sql插入数据方式

	sql: 原生sql
	autoId: 是否是自增主键 true: 自增主键,返回值将返回插入后的自增主键值 / false: 不是自增主键,返回值中的自增主键值为0
	values... : 可变参数, 用于替换 `?` 占位符

return:

	返回值1: 影响的行数
	返回值2: 插入后的自增主键值 , 受autoId参数影响
*/
func (qw *UpdateWrapper) Insert4SQL(autoId bool, sql string, values ...interface{}) (int64, int64) {

	var r orm.RawSeter

	if !qw.txFlag || qw.txOrmer == nil {
		if o == nil {
			return 0, 0
		}
		r = o.Raw(sql, values...)
	} else {
		r = qw.txOrmer.Raw(sql, values...)
	}
	if r == nil {
		logger.Error("RawSeter is null")
		return 0, 0
	}
	res, err := r.Exec()
	if err != nil {
		logger.Error("Insert4SQL failed: %v", err)
		return 0, 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Insert4SQL RowsAffected failed: %v", err)
		return 0, 0
	}
	if autoId {
		lastId, err1 := res.LastInsertId()
		if err1 != nil {
			logger.Error("Get LastInsertId failed: %v", err1)
			return 0, 0
		}
		return count, lastId
	}
	return count, 0
}

// 无事务的自定义插入
func Insert4SQL(autoId bool, sql string, values ...interface{}) (int64, int64) {

	var r orm.RawSeter

	if o == nil {
		return 0, 0
	}
	r = o.Raw(sql, values...)

	if r == nil {
		logger.Error("RawSeter is null")
		return 0, 0
	}
	res, err := r.Exec()
	if err != nil {
		logger.Error("Insert4SQL failed: %v", err)
		return 0, 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("Insert4SQL RowsAffected failed: %v", err)
		return 0, 0
	}
	if autoId {
		lastId, err1 := res.LastInsertId()
		if err1 != nil {
			logger.Error("Get LastInsertId failed: %v", err1)
			return 0, 0
		}
		return count, lastId
	}
	return count, 0
}

/*
原生sql的方式进行更新
return:  影响的行数
*/
func (qw *UpdateWrapper) Update4SQL(sql string, values ...interface{}) int64 {
	return qw.exec4SQL(sql, values...)
}

/*
原生sql的方式进行删除
return: 影响的行数
*/
func (qw *UpdateWrapper) Delete4SQL(sql string, values ...interface{}) int64 {
	return qw.exec4SQL(sql, values...)
}

func (qw *UpdateWrapper) exec4SQL(sql string, values ...interface{}) int64 {

	var r orm.RawSeter

	if !qw.txFlag || qw.txOrmer == nil {
		if o == nil {
			return 0
		}
		r = o.Raw(sql, values...)
	} else {
		r = qw.txOrmer.Raw(sql, values...)
	}
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}

	res, err := r.Exec()
	if err != nil {
		logger.Error("exec4SQL failed: %v", err)
		return 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("exec4SQL RowsAffected failed: %v", err)
		return 0
	}
	return count
}

/*
原生sql的方式进行更新
return:  影响的行数
*/
func Update4SQL(sql string, values ...interface{}) int64 {
	return exec4SQL(sql, values...)
}

/*
原生sql的方式进行删除
return: 影响的行数
*/
func Delete4SQL(sql string, values ...interface{}) int64 {
	return exec4SQL(sql, values...)
}

func exec4SQL(sql string, values ...interface{}) int64 {

	var r orm.RawSeter

	if o == nil {
		return 0
	}
	r = o.Raw(sql, values...)
	if r == nil {
		logger.Error("RawSeter is null")
		return 0
	}

	res, err := r.Exec()
	if err != nil {
		logger.Error("exec4SQL failed: %v", err)
		return 0
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("exec4SQL RowsAffected failed: %v", err)
		return 0
	}
	return count
}
