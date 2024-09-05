package mapper

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-orm/logger"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"strings"
)

type MapperUtils struct {
}

var (
	tableIds        []tableCacheVO
	tableNames      []tableCacheVO
	fieldNames      []fieldCacheVO // 驼峰字段名称
	fieldTableNames []fieldCacheVO // 数据库原始字段名称, 带下划线的
	MySQL           = "mysql"
	Sqlite          = "sqlite3"
	active_db       = ""
	active_log      = false
	o               orm.Ormer
	dataSource      string
	params1         = 50
	params2         = 100
	dbActiveFlag    = false
	transaction     = false // 全局事务
)

func Register(activeDB, connStr string, params ...int) {
	if dbActiveFlag {
		panic("Database is already registered")
	}

	if isEmpty(activeDB) {
		panic("ActiveDB is empty")
	}
	if activeDB != MySQL && activeDB != Sqlite {
		panic(activeDB + " is not support")
	}
	if isEmpty(connStr) {
		panic("Connection string is empty")
	}
	active_db = activeDB
	dataSource = connStr
	for i, v := range params {
		switch i {
		case 0:
			params1 = v
			break
		case 1:
			params2 = v
			break
		}
	}

	if activeDB == MySQL {
		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase("default", "mysql", dataSource)
	} else {
		orm.RegisterDriver("sqlite3", orm.DRSqlite)
		orm.RegisterDataBase("default", "sqlite3", dataSource)
	}

	orm.SetMaxIdleConns("default", params1)
	orm.SetMaxOpenConns("default", params2)

	o = orm.NewOrm()
	dbActiveFlag = true
}

// 开启全局事务 这样的话 BuilderUpdateWrapper 函数就不需要填写第二个参数了
// 如果你填了 BuilderUpdateWrapper 的第二个参数 , 那将以你填的为准
func OpenTransaction() {
	transaction = true
}

func OpenLog() {
	active_log = true
}

// 建议使用该函数去获取ORM 而不是通过NewOrm的方式
func GetOrm() orm.Ormer {
	return o
}

func (qw *UpdateWrapper) removeFalseUpdates() *UpdateWrapper {
	if qw.updates == nil || len(qw.updates) == 0 {
		return qw
	}
	var updateset []updateSet
	for i := 0; i < len(qw.updates); i++ {
		if qw.updates[i].condition {
			updateset = append(updateset, qw.updates[i])
		}
	}
	qw.updates = updateset
	return qw
}

// 查询sql组合器
func (qw *QueryWrapper) queryWrapper4SQL() (string, []interface{}) {
	baseSQL, values := getQuerySQL(qw.query)
	if qw.sorts != nil || len(qw.sorts) > 0 {
		sorts := qw.sorts
		for i := 0; i < len(sorts); i++ {
			sort := sorts[i]
			if i == 0 {
				if sort.condition {
					if sort.isAsc {
						baseSQL = fmt.Sprintf("%s order by %s asc ", baseSQL, sort.columns)
					} else {
						baseSQL = fmt.Sprintf("%s order by %s desc ", baseSQL, sort.columns)
					}
				}
			} else {
				if sort.condition {
					if sort.isAsc {
						baseSQL = fmt.Sprintf("%s , %s asc ", baseSQL, sort.columns)
					} else {
						baseSQL = fmt.Sprintf("%s , %s desc ", baseSQL, sort.columns)
					}
				}
			}
		}
	}

	if !isEmpty(qw.lastSQL) {
		baseSQL = fmt.Sprintf("%s %s", baseSQL, qw.lastSQL)
	}

	return baseSQL, values
}

func getQuerySQL(querys []queryCriteria) (string, []interface{}) {
	baseSQL := ""
	values := make([]interface{}, 0)
	if querys == nil {
		return baseSQL, values
	}
	if len(querys) == 0 {
		return baseSQL, values
	}
	for _, query := range querys {
		if !query.condition {
			continue
		}
		if query.actions == "BETWEEN" || query.actions == "NOT_IN" || query.actions == "IN" {
			if query.actions == "BETWEEN" {
				baseSQL = fmt.Sprintf("%s and %s BETWEEN ? AND ? ", baseSQL, query.columns)
			} else {
				str := "("
				for i := 0; i < len(query.values); i++ {
					if i == len(query.values)-1 {
						str += "?"
					} else {
						str += "?,"
					}
				}
				str += ")"

				baseSQL = fmt.Sprintf("%s and %s %s %s ", baseSQL, query.columns, getSqlKeyword(query.actions), str)
			}
		} else {
			if query.actions == "IS_NULL" || query.actions == "IS_NOT_NULL" {
				baseSQL = fmt.Sprintf("%s and %s %s ", baseSQL, query.columns, getSqlKeyword(query.actions))
			} else {
				baseSQL = fmt.Sprintf("%s and %s %s ? ", baseSQL, query.columns, getSqlKeyword(query.actions))
			}

		}
		for i := 0; i < len(query.values); i++ {
			values = append(values, query.values[i])
		}
	}
	return baseSQL, values
}

func getSqlKeyword(sqlKeyword string) string {
	switch sqlKeyword {
	case "EQ":
		return "="
	case "NE":
		return "<>"
	case "GT":
		return ">"
	case "GE":
		return ">="
	case "LT":
		return "<"
	case "LE":
		return "<="
	case "IN":
		return "in"
	case "BETWEEN":
		return "BETWEEN"
	case "LIKE":
		return "LIKE"
	case "IS_NULL":
		return "IS NULL"
	case "IS_NOT_NULL":
		return "IS NOT NULL"
	case "NOT_IN":
		return "not in"
	case "NOT_LIKE":
		return "not like"
	}
	panic("Unknown sqlKeyword appears")
}

func getTableId(T any) string {
	v := reflect.Indirect(reflect.ValueOf(T))
	//加入缓存 , 速度更快
	tp := v.Type().String()
	idFieldName := getCacheTableId(tp)
	if idFieldName != "null" {
		return idFieldName
	}
	switch v.Kind() {
	case reflect.Slice:
		if v.Len() == 0 {
			return "null"
		}
		return getTableId4Array(T, tp)
	case reflect.Struct:
		return getTableId4Object(T, tp)
	default:
		fmt.Println("Input is neither a slice nor a struct")
	}
	return "null"
}

func getParmaStruct(T any) string {
	v := reflect.Indirect(reflect.ValueOf(T))
	switch v.Kind() {
	case reflect.Slice:
		return "Array"
	case reflect.Struct:
		return "Object"
	default:
		fmt.Println("Input is neither a slice nor a struct")
	}
	return "null"
}

func getTableId4Object(T any, typeName string) string {
	ref := reflect.Indirect(reflect.ValueOf(T)).Type()

	for i := 0; i < ref.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := ref.Field(i)
		// 通过字段名, 找到字段类型信息
		if tp, ok := ref.FieldByName(fieldType.Name); ok {
			// 从tag中取出需要的tag
			if !isEmpty(tp.Tag.Get("tableId")) {
				idName := camelToUnderscore(tp.Tag.Get("tableId"))
				addTableCache(typeName, idName, "basedao-tableIds")
				return idName
			}
		}
	}
	return "null"
}

// 获取排除后的字段列表
/*
	T : 对象
	tableName : 表名
	excludeFields : 排除字段
	autoId : 是否是自增
	idFieldName : id字段名
	excludeEmpty : 是否排除空值
		true: 将排除所有空值字段与数据
		false: 不排除所有空值字段与数据
*/
func getExcludeFiledName(T any, tableName string, excludeFields []string, idFieldName string, excludeEmpty bool) ([]string, []interface{}, bool) {
	// 获取所有的数据库表字段名和传入的对象值
	fields, values := getAllFiledName(T, tableName)
	id := camelToUnderscore(idFieldName)
	names := make([]string, 0)
	vals := make([]interface{}, 0)

	// 根据用户是否输入了ID来判断是否自增
	autoId := hasEmptyId(fields, values, id)

	for i := 0; i < len(fields); i++ {
		field := fields[i]
		val := values[i]
		// 如果该字段是id , 并且为自增,则直接排除
		if field == id && autoId {
			continue
		}
		// 如果是排除项 ,则直接排除
		if isExclude(field, excludeFields) {
			continue
		}
		// 判断是否是空值
		if excludeEmpty && hasEmptyVal(val) {
			continue
		}
		names = append(names, field)
		vals = append(vals, val)
	}

	return names, vals, autoId
}

/*
主键是否为空

	true: 空
	false: 非空
*/
func hasEmptyId(fields []string, values []interface{}, idName string) bool {
	var val interface{}
	for i := 0; i < len(fields); i++ {
		// 如果该字段是id
		if fields[i] == idName {
			val = values[i]
			break
		}
	}
	if val == nil {
		return true
	}
	return hasEmptyVal(val)
}

func hasEmptyVal(val interface{}) bool {
	switch t := val.(type) {
	case string:
		if val.(string) == "" {
			return true
		}
	case int:
		if val.(int) == 0 {
			return true
		}
	case int8:
		if val.(int8) == 0 {
			return true
		}
	case int16:
		if val.(int16) == 0 {
			return true
		}
	case int32:
		if val.(int32) == 0 {
			return true
		}
	case int64:
		if val.(int64) == 0 {
			return true
		}
	case float32:
		if val.(float32) == 0 {
			return true
		}
	case float64:
		if val.(float64) == 0 {
			return true
		}
	default:
		fmt.Printf("The value is of an unknown type: %T\n", t)
	}
	return false
}

func isExclude(field string, excludeFields []string) bool {
	if excludeFields == nil || len(excludeFields) == 0 {
		return false
	}
	for _, excludeField := range excludeFields {
		if excludeField == field {
			return true
		}
	}
	return false
}

func getAllFiledName(T any, tableName string) ([]string, []interface{}) {

	fn := getCatchFieldNames(tableName)
	ftn := getCatchFieldTableNames(tableName)

	if fn == nil {
		names := make([]string, 0)
		filedNameList := make([]string, 0)
		ref := reflect.Indirect(reflect.ValueOf(T)).Type()
		for i := 0; i < ref.NumField(); i++ {
			field := ref.Field(i)
			//驼峰转下划线方式
			name := camelToUnderscore(field.Name)
			// 驼峰字段名称
			names = append(names, field.Name)
			// 转成数据库格式的字段名称 (带下划线的)
			filedNameList = append(filedNameList, name)
		}
		fieldNames = append(fieldNames, fieldCacheVO{Name: tableName, Fields: names})
		fieldTableNames = append(fieldTableNames, fieldCacheVO{Name: tableName, Fields: filedNameList})
		fn = names
		ftn = filedNameList
	}
	values := make([]interface{}, 0)
	obj := reflect.Indirect(reflect.ValueOf(T))

	for i := 0; i < len(fn); i++ {
		val := obj.FieldByName(fn[i]).Interface()
		values = append(values, val)
	}
	return ftn, values
}

// 从缓存中取原始格式的字段名
func getCatchFieldNames(tableName string) []string {
	if fieldNames != nil && len(fieldNames) > 0 {
		for _, vo := range fieldNames {
			if vo.Name == tableName {
				return vo.Fields
			}
		}
	}
	return nil
}

// 从缓存中取数据库格式的字段名
func getCatchFieldTableNames(tableName string) []string {
	if fieldTableNames != nil && len(fieldTableNames) > 0 {
		for _, vo := range fieldTableNames {
			if vo.Name == tableName {
				return vo.Fields
			}
		}
	}
	return nil
}

func getTableId4Array(T any, typeName string) string {
	v := reflect.Indirect(reflect.ValueOf(T))
	elemType := reflect.Indirect(v.Index(0)).Type()
	for i := 0; i < elemType.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := elemType.Field(i)
		// 通过字段名, 找到字段类型信息
		if tp, ok := elemType.FieldByName(fieldType.Name); ok {
			// 从tag中取出需要的tag
			if !isEmpty(tp.Tag.Get("tableId")) {
				idName := camelToUnderscore(tp.Tag.Get("tableId"))
				addTableCache(typeName, idName, "basedao-tableIds")
				return idName
			}
		}
	}
	return "null"
}

func getTableName(T interface{}) string {
	funcName := "TableName"
	tp := getParmaStruct(T)
	if tp == "null" {
		return "null"
	}

	typeName := reflect.ValueOf(T).Type().String()
	tableName := getCacheTableName(typeName)
	if tableName != "null" {
		return tableName
	}

	if tp == "Object" {
		methodValue := reflect.ValueOf(T).MethodByName(funcName)
		if !methodValue.IsValid() {
			return "null"
		}
		back := methodValue.Call(nil)
		tbName := back[0].String()
		addTableCache(typeName, tbName, "basedao-tableNames")
		return tbName
	} else {
		v := reflect.Indirect(reflect.ValueOf(T))
		elemType := v.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		obj := reflect.New(elemType).Elem().Interface()
		funcVal := reflect.ValueOf(obj).MethodByName(funcName)
		// methodValue, _ := elemType.MethodByName(funcName)
		// fmt.Println(methodValue.Name)
		// funcVal := methodValue.Func

		if !funcVal.IsValid() {
			return "null"
		}
		back := funcVal.Call(nil)
		tbName := back[0].String()
		addTableCache(typeName, tbName, "basedao-tableNames")
		return tbName
	}

}

func isEmpty(str string) bool {
	return len(str) == 0
}

func camelToUnderscore(str string) string {

	var result strings.Builder
	// 将字符串首字母转换为小写
	str = strings.ToLower(str[0:1]) + str[1:]

	for i, char := range str {
		if char >= 'A' && char <= 'Z' {
			// 如果当前字符是大写字母，则在前面添加下划线，并将该字母转为小写
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteByte(byte(char + 32))
		} else {
			result.WriteByte(byte(char))
		}
	}

	return result.String()
}

func getCacheTableId(tp string) string {
	if isEmpty(tp) {
		return "null"
	}
	// tableIds := utils.GetGlobalCache("basedao-tableIds")
	// if tableIds == nil {
	// 	return "null"
	// }
	// vos := tableIds.([]tableCacheVO)
	if len(tableIds) == 0 {
		return "null"
	}
	for i := 0; i < len(tableIds); i++ {
		if tableIds[i].ObjType == tp {
			// fmt.Println("从缓存中取得TableId成功", tableIds[i].Name)
			return tableIds[i].Name
		}
	}

	return "null"
}
func getCacheTableName(tp string) string {
	if isEmpty(tp) {
		return "null"
	}
	// tableIds := utils.GetGlobalCache("basedao-tableNames")
	// if tableIds == nil {
	// 	return "null"
	// }
	// vos := tableIds.([]tableCacheVO)
	if len(tableNames) == 0 {
		return "null"
	}
	for i := 0; i < len(tableNames); i++ {
		if tableNames[i].ObjType == tp {
			// fmt.Println("从缓存中取得TableName成功", tableNames[i].Name)
			return tableNames[i].Name
		}
	}

	return "null"
}

func addTableCache(tp string, name string, cacheName string) {
	// tableIds := utils.GetGlobalCache("tableIds")

	if cacheName == "basedao-tableNames" {
		if len(tableNames) == 0 {
			// var vos []tableCacheVO
			// vos = append(vos, tableCacheVO{ObjType: tp, Name: name})
			// utils.SetGlobalCache(cacheName, vos)
			tableNames = append(tableNames, tableCacheVO{ObjType: tp, Name: name})
			return
		}
		for i := 0; i < len(tableNames); i++ {
			if tableNames[i].ObjType == tp {
				return
			}
		}
		tableNames = append(tableNames, tableCacheVO{ObjType: tp, Name: name})
		return
		// utils.SetGlobalCache(cacheName, vos)
	}

	if len(tableIds) == 0 {
		// var vos []tableCacheVO
		// vos = append(vos, tableCacheVO{ObjType: tp, Name: name})
		// utils.SetGlobalCache(cacheName, vos)
		tableIds = append(tableIds, tableCacheVO{ObjType: tp, Name: name})
		return
	}
	for i := 0; i < len(tableIds); i++ {
		if tableIds[i].ObjType == tp {
			return
		}
	}
	tableIds = append(tableIds, tableCacheVO{ObjType: tp, Name: name})

	// if tableIds == nil {
	// 	var vos []tableCacheVO
	// 	vos = append(vos, tableCacheVO{ObjType: tp, Name: name})
	// 	utils.SetGlobalCache(cacheName, vos)
	// 	return
	// }
	// vos := tableIds.([]tableCacheVO)
	// if len(vos) == 0 {
	// 	vos = append(vos, tableCacheVO{ObjType: tp, Name: name})
	// 	utils.SetGlobalCache(cacheName, vos)
	// 	return
	// }
	// for i := 0; i < len(vos); i++ {
	// 	if vos[i].ObjType == tp {
	// 		return
	// 	}
	// }
	// vos = append(vos, tableCacheVO{ObjType: tp, Name: name})
	// utils.SetGlobalCache(cacheName, vos)
	// return
}

func checkErr(err error) {
	if err != nil {
		logger.Info("database error: %v", err)
	}
}

func saveLastInsertId(T interface{}, lastId int64) {
	ref := reflect.ValueOf(T).Elem().Type()
	idFieldName := ""
	for i := 0; i < ref.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		field := ref.Field(i)
		// 通过字段名, 找到字段类型信息
		if tp, ok := ref.FieldByName(field.Name); ok {
			// 从tag中取出需要的tag
			if !isEmpty(tp.Tag.Get("tableId")) {
				idFieldName = field.Name
			}
		}
	}
	obj := reflect.ValueOf(T).Elem()
	idField := obj.FieldByName(idFieldName)
	if idField.IsValid() && idField.CanSet() {
		switch idField.Type().Kind() {
		case reflect.Int32:
			idField.SetInt(lastId)
		case reflect.Int64:
			idField.SetInt(lastId)
		case reflect.Int16:
			idField.SetInt(lastId)
		case reflect.Int8:
			idField.SetInt(lastId)
		case reflect.Int:
			idField.SetInt(lastId)
		}
	}
}

func LogInfo(methodName string, format string) {
	if active_log {
		logger.Info("[%s]: ==> %s", methodName, format)
	}
}

func isSlice(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	kind := value.Kind()

	return kind == reflect.Slice
}
