package mapper

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"gorm.io/gorm"
)

type BaseMapper struct {
}

func Page[T any](page *PageUtils) PageData[T] {
	qw := page.wrapper

	offset := page.getOffSet()
	pageSize := page.getPageSize()

	qw.lastSQL = fmt.Sprintf("limit %d,%d", offset, pageSize)
	list := SelectList[T](qw)

	qw.lastSQL = ""
	count := SelectCount[T](qw)

	page.setTotalSize(int64(count))

	return buildPageData(page, list)
}

// ConvertPageData 将 PageData[T] 转换为 PageData[T2]
func ConvertPageData[T any, T2 any](pageData PageData[T], list []T2) PageData[T2] {
	return pagePo2Vo(pageData, list)
}

/**
 * 根据主键查询
 * 该查询会根据传入的参数进行两种结果类型封装
 * 1. 当传入为对象类型时, 会将结果集封装为一个对象(不推荐,因为go的空值没有null这个结果,只能判断id是否为0)
 * 2. 当传入对象为数组类型时 , 会将结果集封装为一个数组(推荐,可以通过判断数组length是否为0来判断是否存在数据)
 */
func SelectById[T any](id interface{}) []T {
	var zero T
	idFieldName := getTableId(zero)
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}

	tableName := getTableName(zero)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	tp := getParmaStruct(zero)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}

	sql := fmt.Sprintf("select * from %s where %s = ?", tableName, idFieldName)
	LogInfo("SelectById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("SelectById", fmt.Sprintf("Parameters: %v", id))

	result, err := gorm.G[T](globalDB).Raw(sql, id).Find(context.Background())
	if err != nil {
		logger.Error("SelectList failed: {}", err)
		return result
	}
	return result
}

func SelectList[T any](qw *QueryWrapper) []T {

	if qw == nil {
		logger.Error("SelectList failed: QueryWrapper is nil")
		return nil
	}

	var zero T
	tableName := getTableName(zero)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	baseSQL, values := qw.queryWrapper4SQL()

	baseSQL = fmt.Sprintf("select * from %s where 1=1 %s", tableName, baseSQL)
	LogInfo("SelectList", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("SelectList", fmt.Sprintf("Parameters: %v", values))

	result, err := gorm.G[T](globalDB).Raw(baseSQL, values...).Find(context.Background())
	if err != nil {
		logger.Error("SelectList failed: {}", err)
		return result
	}
	return result
}

func SelectCount[T any](qw *QueryWrapper) int64 {

	if qw == nil {
		qw = BuilderQueryWrapper()
	}

	var zero T
	tableName := getTableName(zero)
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

	count, err := gorm.G[int](globalDB).Raw(baseSQL, values...).Find(context.Background())
	if err != nil {
		logger.Error("SelectCount failed: {}", err)
		return 0
	}
	if len(count) == 0 {
		return 0
	}
	return int64(count[0])
}

func Delete[T any](qw *UpdateWrapper) int64 {

	if qw == nil {
		logger.Error("Delete failed: UpdateWrapper is nil")
		return 0
	}

	var zero T
	tableName := getTableName(zero)
	if tableName == "null" {
		panic("TableName method does not exist")
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

	db := globalDB
	if qw.txFlag && qw.txOrmer != nil {
		db = qw.txOrmer
	}
	result := gorm.WithResult()
	// Execute with parameters
	err := gorm.G[any](db, result).Exec(context.Background(), baseSQL, values...)

	if err != nil {
		logger.Error("Delete failed: {}", err)
		return 0
	}
	return result.RowsAffected
}

func DeleteByIds[T any](ids interface{}, qw *UpdateWrapper) int64 {

	if qw == nil {
		qw = BuilderUpdateWrapper(false)
	}

	var zero T
	idFieldName := getTableId(zero)
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}

	tableName := getTableName(zero)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	db := globalDB
	if qw.txFlag && qw.txOrmer != nil {
		db = qw.txOrmer
	}
	// 获取 ids 的反射值
	idsVal := reflect.ValueOf(ids)
	vals := make([]interface{}, idsVal.Len())
	var str strings.Builder
	str.WriteString("(")
	// 遍历切片
	for i := 0; i < idsVal.Len(); i++ {
		vals[i] = idsVal.Index(i).Interface()
		if i == idsVal.Len()-1 {
			str.WriteString("?")
		} else {
			str.WriteString("?,")
		}
	}
	str.WriteString(")")

	sql := fmt.Sprintf("delete from %s where %s in %s", tableName, idFieldName, str.String())
	LogInfo("DeleteByIds", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("DeleteByIds", fmt.Sprintf("Parameters: %v", vals))
	result := gorm.WithResult()
	// Execute with parameters
	err := gorm.G[any](db, result).Exec(context.Background(), sql, vals...)

	if err != nil {
		logger.Error("Delete failed: {}", err)
		return 0
	}
	return result.RowsAffected
}

// 使用事务的删除
func DeleteById[T any](id interface{}, qw *UpdateWrapper) int64 {

	if qw == nil {
		qw = BuilderUpdateWrapper(false)
	}

	var zero T
	idFieldName := getTableId(zero)
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}

	tableName := getTableName(zero)
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	db := globalDB
	if qw.txFlag && qw.txOrmer != nil {
		db = qw.txOrmer
	}
	sql := fmt.Sprintf("delete from %s where %s = ?", tableName, idFieldName)
	LogInfo("DeleteById", fmt.Sprintf("Preparing: %s", sql))
	LogInfo("DeleteById", fmt.Sprintf("Parameters: %v", id))

	result := gorm.WithResult()
	// Execute with parameters
	err := gorm.G[any](db, result).Exec(context.Background(), sql, id)

	if err != nil {
		logger.Error("Delete failed: {}", err)
		return 0
	}
	return result.RowsAffected
}

func Update[T any](qw *UpdateWrapper) int64 {

	if qw == nil {
		logger.Error("Update failed: UpdateWrapper is nil")
		return 0
	}

	var zero T
	tp := getParmaStruct(zero)
	if tp == "null" {
		panic("Input is neither a slice nor a struct")
	}
	tableName := getTableName(zero)
	if tableName == "null" {
		panic("TableName method does not exist")
	}
	if qw.updates == nil || len(qw.updates) == 0 {
		panic("Missing fields to be updated")
	}
	values := make([]interface{}, 0)
	baseSQL := fmt.Sprintf("update %s ", tableName)

	// 去除flag为false的更新字段
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

	db := globalDB
	if qw.txFlag && qw.txOrmer != nil {
		db = qw.txOrmer
	}

	result := gorm.WithResult()
	// Execute with parameters
	err := gorm.G[any](db, result).Exec(context.Background(), baseSQL, values...)

	if err != nil {
		logger.Error("Update failed: {}", err)
		return 0
	}
	return result.RowsAffected
}

// 批量插入 - 可极大提高效率
// bulk, 单次插入数量
func InsertAll[T any](list []T, qw *UpdateWrapper) int64 {

	if qw == nil {
		qw = BuilderUpdateWrapper(false)
	}

	bulk := qw.bulk

	// 检查 list 是否为 nil 或空切片
	if list == nil || len(list) == 0 {
		return 0
	}

	// 将 []T 转换为 []interface{}
	params := make([]interface{}, len(list))
	for i, v := range list {
		params[i] = v
	}

	idFieldName := getTableId(params[0])
	if idFieldName == "null" {
		panic("Field 'Primary key' does not exist , Please check if the Tag of the primary key attribute in the entity class contains the tableId attribute")
	}
	tableName := getTableName(params[0])
	if tableName == "null" {
		panic("TableName method does not exist")
	}

	return insertMulti(idFieldName, tableName, params, bulk, qw)
}

func insertMulti(idFieldName, tableName string, params []interface{}, bulk int, qw *UpdateWrapper) int64 {

	sind := reflect.ValueOf(params)
	if sind.Len() == 0 {
		return 0
	}

	// 取第一条数据，确定字段结构
	fields, _, _ := getExcludeFiledName(
		sind.Index(0).Interface(),
		tableName,
		qw.excludeField,
		idFieldName,
		qw.excludeEmpty,
	)
	if len(fields) == 0 {
		panic("No fields were found")
	}

	// =========================
	// 1. 构建字段 SQL：( field1, field2, field3 )
	// =========================
	var colBuilder strings.Builder
	colBuilder.WriteString("(")
	for i, f := range fields {
		colBuilder.WriteString(" ")
		if i > 0 {
			colBuilder.WriteString(", ")
		}
		colBuilder.WriteString(f)
	}
	colBuilder.WriteString(" ) ")

	sql1 := colBuilder.String()

	// =========================
	// 2. baseSQL
	// =========================
	baseSQL := "insert into " + tableName + " " + sql1 + " values "

	// =========================
	// 3. 生成单行占位符模板：(?, ?, ?)
	// =========================
	rowTpl := buildRowPlaceholder(len(fields))

	// =========================
	// 4. 计算批次数
	// =========================
	if bulk <= 0 {
		bulk = sind.Len()
	}

	endOffset := (sind.Len() + bulk - 1) / bulk

	var sum int64 = 0

	// =========================
	// 5. 分批插入
	// =========================
	for offset := 0; offset < endOffset; offset++ {
		startIndex := offset * bulk
		endIndex := startIndex + bulk
		if endIndex > sind.Len() {
			endIndex = sind.Len()
		}

		rows := endIndex - startIndex
		if rows <= 0 {
			continue
		}

		// 5.1 构建多行占位符
		valuesSQL := buildMultiRowPlaceholder(rowTpl, rows)

		// 5.2 收集 values
		values := make([]interface{}, 0, rows*len(fields))
		for i := startIndex; i < endIndex; i++ {
			_, vals, _ := getExcludeFiledName(
				sind.Index(i).Interface(),
				tableName,
				qw.excludeField,
				idFieldName,
				qw.excludeEmpty,
			)
			values = append(values, vals...)
		}

		// =========================
		// 6. 执行 SQL（事务）
		// =========================
		db := globalDB.Begin()
		result := gorm.WithResult()

		err := gorm.G[any](db, result).
			Exec(context.Background(), baseSQL+valuesSQL, values...)
		LogInfo("InsertAll", fmt.Sprintf("Preparing: %v", baseSQL+valuesSQL))
		LogInfo("InsertAll", fmt.Sprintf("Parameters: %v", values))

		if err != nil {
			logger.Error("Insert failed: {}", err)
			db.Rollback()
			return 0
		}

		db.Commit()
		sum += result.RowsAffected
	}

	return sum
}

func buildMultiRowPlaceholder(rowTpl string, rows int) string {
	if rows <= 0 {
		return ""
	}

	var b strings.Builder
	// 预估容量，避免多次扩容
	b.Grow(len(rowTpl)*rows + (rows-1)*2)

	for i := 0; i < rows; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(rowTpl)
	}

	return b.String()
}

func buildRowPlaceholder(fieldCount int) string {
	var b strings.Builder
	b.Grow(fieldCount * 3)

	b.WriteString("(")
	for i := 0; i < fieldCount; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("?")
	}
	b.WriteString(")")

	return b.String()
}

func Insert[T any](pojo *T, qw *UpdateWrapper) int64 {
	if qw == nil {
		qw = BuilderUpdateWrapper(false)
	}

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

	excludeField := make([]string, 0)
	excludeEmpty := false
	if qw != nil {
		excludeField = qw.excludeField
		excludeEmpty = qw.excludeEmpty
	}

	fields, values, autoId := getExcludeFiledName(pojo, tableName, excludeField, idFieldName, excludeEmpty)
	if len(fields) == 0 {
		panic("No fields were found")
	}

	// 使用 strings.Builder 优化字符串拼接
	var sb strings.Builder
	sb.WriteString("insert into ")
	sb.WriteString(tableName)
	sb.WriteString(" (")

	for i, field := range fields {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(field)
	}
	sb.WriteString(") values (")

	for i := 0; i < len(fields); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("?")
	}
	sb.WriteString(")")

	baseSQL := sb.String()

	LogInfo("Insert", fmt.Sprintf("Preparing: %s", baseSQL))
	LogInfo("Insert", fmt.Sprintf("Parameters: %v", values))

	db := globalDB
	if qw != nil && qw.txFlag && qw.txOrmer != nil {
		db = qw.txOrmer
	}
	result := gorm.WithResult()
	// Execute with parameters
	err := gorm.G[any](db, result).Exec(context.Background(), baseSQL, values...)

	if err != nil {
		logger.Error("Update failed: {}", err)
		return 0
	}

	if autoId {
		lastId, err1 := result.Result.LastInsertId()
		if err1 != nil {
			logger.Error("Get LastInsertId failed: {}", err1)
			return 0
		}
		saveLastInsertId(pojo, lastId)
	}
	return result.RowsAffected
}

func SelectCount4SQL(sql string, values ...interface{}) int64 {
	count, err := gorm.G[int](globalDB).Raw(sql, values...).Find(context.Background())
	if err != nil {
		logger.Error("SelectCount4SQL failed: {}", err)
		return 0
	}
	if len(count) == 0 {
		return 0
	}
	return int64(count[0])
}
func SelectList4SQL[T any](sql string, values ...interface{}) []T {
	result, err := gorm.G[T](globalDB).Raw(sql, values...).Find(context.Background())
	if err != nil {
		logger.Error("SelectList4SQL failed: {}", err)
		return result
	}
	return result
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
func Insert4SQL(qw *UpdateWrapper) (int64, int64) {
	if qw == nil {
		logger.Error("Insert4SQL failed: UpdateWrapper is nil")
		return 0, 0
	}
	if isEmpty(qw.baseSql) {
		logger.Error("Insert4SQL failed: baseSql is null")
		return 0, 0
	}

	db := globalDB
	if qw.txFlag && qw.txOrmer != nil {
		db = qw.txOrmer
	}

	// Execute raw SQL
	result := gorm.WithResult()
	// Execute with parameters
	err := gorm.G[any](db, result).Exec(context.Background(), qw.baseSql, qw.values...)

	if err != nil {
		logger.Error("Insert4SQL failed: {}", err)
		return 0, 0
	}
	count := result.RowsAffected
	if qw.autoId {
		lastId, err1 := result.Result.LastInsertId()
		if err1 != nil {
			logger.Error("Get LastInsertId failed: {}", err1)
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
func Update4SQL(qw *UpdateWrapper) int64 {
	if qw == nil {
		logger.Error("Update4SQL failed: UpdateWrapper is nil")
		return 0
	}
	if isEmpty(qw.baseSql) {
		logger.Error("Update4SQL failed: baseSql is null")
		return 0
	}
	return exec4SQL(qw)
}

/*
原生sql的方式进行删除
return: 影响的行数
*/
func Delete4SQL(qw *UpdateWrapper) int64 {
	if qw == nil {
		logger.Error("Delete4SQL failed: UpdateWrapper is nil")
		return 0
	}
	if isEmpty(qw.baseSql) {
		logger.Error("Delete4SQL failed: baseSql is null")
		return 0
	}
	return exec4SQL(qw)
}

func exec4SQL(qw *UpdateWrapper) int64 {

	db := globalDB
	if qw.txFlag && qw.txOrmer != nil {
		db = qw.txOrmer
	}
	// Execute raw SQL
	result := gorm.WithResult()
	// Execute with parameters
	err := gorm.G[any](db, result).Exec(context.Background(), qw.baseSql, qw.values...)
	if err != nil {
		logger.Error("EXEC SQL failed: {}", err)
		return 0
	}
	return result.RowsAffected
}
