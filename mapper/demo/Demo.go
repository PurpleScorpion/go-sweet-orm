package mapperDemo

import (
	"encoding/json"
	"fmt"
	"go-sweet-orm/mapper"
)

type Demo struct {
}

// 根据id查询
func demo1() {
	var log []Logs
	mapper.SelectById(&log, 1)
	fmt.Println(log)
}

// 查询列表
func demo2() {
	var log_type = "設定変更"

	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	// 参数1：条件是否生效 , 参数2: 数据库列名 , 参数3: 条件值
	qw.Eq(isEmpty(log_type), "log_type", log_type)
	// 与上方一致 , 需注意在Like下参数3是不需要写 % 的
	qw.Like(true, "log_content", "逆潮流防止")
	// 参数1：条件是否生效 , 参数2: 数据库列名 , 参数3: 条件值(可变参数,可写多个)
	qw.In(true, "log_level", "ユーザー", "システム")
	// 若排序是时间格式 , 可以使用 OrderByTimeAsc
	qw.OrderByTimeAsc(true, "log_time")
	// qw.OrderByAsc(true, "log_time")
	mapper.SelectList(qw)
	fmt.Println(log)
}

// 查询列表 - 全部 此功能慎用 , 查询到结果过多时可能导致服务器崩溃
func demo2_all() {
	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	mapper.SelectList(qw)
	fmt.Println(log)
}

// 查询count
func demo3() {
	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	qw.Like(true, "log_content", "逆潮流防止")

	//注意 , count查询下该OrderBy条件是不生效的 , 也就是你写不写无所谓
	qw.OrderByAsc(true, "log_time")
	count := mapper.SelectCount(qw)
	fmt.Println(fmt.Sprintf("查询到的count: %s", count))
}

// 查询count - 全部
func demo3_all() {
	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	count := mapper.SelectCount(qw)
	fmt.Println(fmt.Sprintf("查询到的count: %s", count))
}

// 原始sql语句查询列表 , 此项不检查主键与表名 , 可用vo类等当做结果集
func demo4() {
	var log []Logs
	mapper.SelectList4SQL(&log, "select * from logs where log_type = ?", "設定変更")
	fmt.Println(log)
}

// 原始sql语句查询数量
func demo5() {
	count := mapper.SelectCount4SQL("select count(*) from logs where log_type = ?", "設定変更")
	fmt.Println(fmt.Sprintf("查询到的count: %s", count))
}

// 更新
func demo6() {
	var log Logs
	qw := mapper.BuilderQueryWrapper(&log)
	// Set函数可以调用多次, 但不可以一次都不调用
	qw.Set(true, "log_type", "設定変更")
	qw.Set(true, "log_level", "INFO")
	qw.Eq(true, "id", 1)
	count := mapper.Update(qw)
	fmt.Println(fmt.Sprintf("影响行数: %s", count))
}

// 更新 - 全部 此功能慎用 会导致全表所有数据更新
func demo6_all() {
	var log Logs
	qw := mapper.BuilderQueryWrapper(&log)
	qw.Set(true, "log_type", "設定変更")
	qw.Set(true, "log_level", "INFO")
	count := mapper.Update(qw)
	fmt.Println(fmt.Sprintf("影响行数: %s", count))
}

// 根据ID删除数据
func demo7() {
	var log Logs
	count := mapper.DeleteById(&log, 1)
	fmt.Println(fmt.Sprintf("影响行数: %s", count))
}

// 根据条件删除
func demo8() {
	var log Logs
	qw := mapper.BuilderQueryWrapper(&log)
	qw.Eq(true, "log_type", "設定変更")
	count := mapper.Delete(qw)
	fmt.Println(fmt.Sprintf("影响行数: %s", count))
}

// 删除全部 此功能慎用 , 会清空表数据
func demo8_all() {
	var log Logs
	qw := mapper.BuilderQueryWrapper(&log)
	count := mapper.Delete(qw)
	fmt.Println(fmt.Sprintf("影响行数: %s", count))
}

// 新增 默认自增主键 空值排除
func demo9() {
	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "ユーザー"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "設定変更"

	mapper.Insert(&resLog)
	// 新增后该对象中会有新增后的id
	fmt.Println(resLog)
}

// 新增 默认自增主键 空值排除 指定排除某些字段
func demo10() {
	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "ユーザー"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "設定変更"
	// 排除的字段为表字段名 , 可写多个 , 写上之后 , 即使上方赋值也不会插入该值
	mapper.Insert(&resLog, "log_time", "created_date")
	// 新增后该对象中会有新增后的id
	fmt.Println(resLog)
}

// 新增自定义主键规则与排除字段
func demo11() {
	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "ユーザー"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "設定変更"
	// 第一个true为主键是否为自增(true为自增) 第二个false为是否排除空值(false为不排除空值)
	mapper.InsertCustom(&resLog, true, false, "log_time", "created_date")
	// 新增后该对象中会有新增后的id
	fmt.Println(resLog)
}

// 分页工具
func demo12() {
	// 必须: 必须是数组类型
	var resLog []Logs
	// 查询器: 必须,可结合条件查询
	qw := mapper.BuilderQueryWrapper(&resLog)
	// 创建分页查询器 必须
	page := mapper.BuilderPageUtils(2, 10, qw)
	// 必须: 必须接收返回值 , 可自行修改 PageData 类中的builder函数和Field来完成符合自己需求的分页结果集
	pageData := mapper.Page(page)
	// 打印结果集 - 测试观察使用
	pd, _ := json.Marshal(pageData)
	fmt.Println(string(pd))

}

func isEmpty(str string) bool {
	return len(str) == 0
}
