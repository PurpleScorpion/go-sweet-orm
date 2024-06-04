package mapper

import (
	"fmt"

	"github.com/beego/beego/orm"
)

type BaseMapper struct {
}

func Page(page PageUtils) PageData {
	qw := page.wrapper
	offSet := page.getOffSet()
	pageSize := page.getPageSize()
	qw.lastSQL = fmt.Sprintf("limit %d,%d", offSet, pageSize)
	SelectList(qw)
	qw.lastSQL = ""
	count := SelectCount(qw)
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
func SelectById(obj interface{}, id interface{}) {
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

	o := orm.NewOrm()

	sql := fmt.Sprintf("select * from %s where %s = ?", tableName, idFieldName)
	LogInfo("SelectById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("SelectById", fmt.Sprintf("Parameters: %v", id))
	r := o.Raw(sql, id)

	if tp == "Array" {
		r.QueryRows(obj)
	} else {
		r.QueryRow(obj)
	}
}

func SelectList(qw QueryWrapper) {
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

	baseSQL, values := queryWrapper4SQL(qw)

	baseSQL = fmt.Sprintf("select * from %s where 1=1 %s", tableName, baseSQL)
	LogInfo("SelectList", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("SelectList", fmt.Sprintf("Parameters: %v", values))
	o := orm.NewOrm()
	r := o.Raw(baseSQL, values...)
	r.QueryRows(resList)
	qw.resList = resList
}

func SelectCount(qw QueryWrapper) int64 {
	resList := qw.resList
	tp := getParmaStruct(resList)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	tableName := getTableName(resList)
	if tableName == "null" {
		panic("TableName method does not exist")
	}
	baseSQL, values := getQuerySQL(qw)
	if !isEmpty(qw.lastSQL) {
		baseSQL = fmt.Sprintf("%s %s", baseSQL, qw.lastSQL)
	}
	baseSQL = fmt.Sprintf("select count(*) from %s where 1=1 %s", tableName, baseSQL)
	LogInfo("SelectCount", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("SelectCount", fmt.Sprintf("Parameters: %v", values))
	var count int64
	o := orm.NewOrm()
	o.Raw(baseSQL, values...).QueryRow(&count)
	return count
}

func Delete(qw QueryWrapper) int64 {
	resList := qw.resList
	tp := getParmaStruct(resList)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	tableName := getTableName(resList)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	o := orm.NewOrm()
	baseSQL := fmt.Sprintf("delete from %s where 1=1 ", tableName)
	// 查询组合器
	sql, values := getQuerySQL(qw)
	// 拼接sql
	baseSQL = fmt.Sprintf("%s %s ", baseSQL, sql)
	// 拼接lastsql
	if !isEmpty(qw.lastSQL) {
		baseSQL = fmt.Sprintf("%s %s", baseSQL, qw.lastSQL)
	}
	LogInfo("Delete", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("Delete", fmt.Sprintf("Parameters: %v", values))
	res, err := o.Raw(baseSQL, values...).Exec()
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

func DeleteById(obj interface{}, id interface{}) int64 {
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
	o := orm.NewOrm()
	sql := fmt.Sprintf("delete from %s where %s = ?", tableName, idFieldName)
	LogInfo("DeleteById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("DeleteById", fmt.Sprintf("Parameters: %v", id))
	res, err := o.Raw(sql, id).Exec()
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

func Update(qw QueryWrapper) int64 {
	resList := qw.resList
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
	sql, vals := getQuerySQL(qw)

	baseSQL = fmt.Sprintf("%s %s ", baseSQL, sql)
	if !isEmpty(qw.lastSQL) {
		baseSQL = fmt.Sprintf("%s %s ", baseSQL, qw.lastSQL)
	}

	for i := 0; i < len(vals); i++ {
		values = append(values, vals[i])
	}
	LogInfo("Update", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("Update", fmt.Sprintf("Parameters: %v", values))
	o := orm.NewOrm()
	res, err := o.Raw(baseSQL, values...).Exec()
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

// 默认 自增主键和空值排除的新增
func Insert(pojo interface{}, excludeField ...string) int64 {
	return baseInsert(pojo, true, true, excludeField...)
}

// 默认自增主键和自定义是否排除空值的新增
func InsertAutoId(pojo interface{}, excludeEmpty bool, excludeField ...string) int64 {
	return baseInsert(pojo, true, excludeEmpty, excludeField...)
}

// 完全自定义的新增
func InsertCustom(pojo interface{}, autoId bool, excludeEmpty bool, excludeField ...string) int64 {
	return baseInsert(pojo, autoId, excludeEmpty, excludeField...)
}

/*
pojo: 实体对象
autoId: 是否为自增主键
excludeEmpty: 是否排除空值
excludeField: 排除字段 - 数据库表字段名
*/
func baseInsert(pojo interface{}, autoId bool, excludeEmpty bool, excludeField ...string) int64 {
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
		panic("Input is neither a slice nor a struct")
	}
	if tp == "slice" {
		panic("The pojo parameter cannot be of array type")
	}

	fields, values := getExcludeFiledName(pojo, tableName, excludeField, autoId, idFieldName, excludeEmpty)
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
	o := orm.NewOrm()
	res, err := o.Raw(baseSQL, values...).Exec()
	checkErr(err)
	if autoId {
		lastId, err := res.LastInsertId()
		checkErr(err)
		saveLastInsertId(pojo, lastId)
	}
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

func SelectCount4SQL(sql string, values ...interface{}) int64 {
	var count int64
	o := orm.NewOrm()
	o.Raw(sql, values...).QueryRow(&count)
	return count
}
func SelectList4SQL(resList interface{}, sql string, values ...interface{}) {
	tp := getParmaStruct(resList)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	if tp != "Array" {
		panic("The result set parameter must be of array type")
	}
	o := orm.NewOrm()
	r := o.Raw(sql, values...)
	r.QueryRows(resList)
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
func Insert4SQL(autoId bool, sql string, values ...interface{}) (int64, int64) {
	o := orm.NewOrm()
	res, err := o.Raw(sql, values...).Exec()
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	if autoId {
		lastId, err := res.LastInsertId()
		checkErr(err)
		return count, lastId
	}
	return count, 0
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
	o := orm.NewOrm()
	res, err := o.Raw(sql, values...).Exec()
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}
