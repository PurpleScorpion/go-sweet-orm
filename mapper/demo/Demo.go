package mapperDemo

import (
	"encoding/json"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-orm/mapper"
)

// import (
//
//	"encoding/json"
//	"fmt"
//	"github.com/PurpleScorpion/go-sweet-orm/mapper"
//
// )
//
// type Demo struct {
// }
//
// 根据id查询
func demo1() {
	var log []Logs
	mapper.SelectById(&log, 1)
	fmt.Println(log)
}

// 查询列表
func demo2() {
	var log_type = "test"

	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	// 参数1：条件是否生效 , 参数2: 数据库列名 , 参数3: 条件值
	qw.Eq(isEmpty(log_type), "log_type", log_type)
	// 与上方一致 , 需注意在Like下参数3是不需要写 % 的
	qw.Like(true, "log_content", "test")
	// 参数1：条件是否生效 , 参数2: 数据库列名 , 参数3: 条件值(可变参数,可写多个)
	qw.In(true, "log_level", "user", "sys")
	// 若排序是时间格式 , 可以使用 OrderByTimeAsc
	qw.OrderByTimeAsc(true, "log_time")
	// qw.OrderByAsc(true, "log_time")
	qw.SelectList()
	fmt.Println(log)
}

// 查询列表 - 全部 此功能慎用 , 查询到结果过多时可能导致服务器崩溃
func demo2_all() {
	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	qw.SelectList()
	fmt.Println(log)
}

// 查询count
func demo3() {
	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	qw.Like(true, "log_content", "test")

	//注意 , count查询下该OrderBy条件是不生效的 , 也就是你写不写无所谓
	qw.OrderByAsc(true, "log_time")
	count := qw.SelectCount()
	fmt.Println(fmt.Sprintf("查询到的count: %d", count))
}

// 查询count - 全部
func demo3_all() {
	var log []Logs
	qw := mapper.BuilderQueryWrapper(&log)
	count := qw.SelectCount()
	fmt.Println(fmt.Sprintf("查询到的count: %d", count))
}

// 原始sql语句查询列表 , 此项不检查主键与表名 , 可用vo类等当做结果集
func demo4() {
	var log []Logs
	mapper.SelectList4SQL(&log, "select * from logs where log_type = ?", "test")
	fmt.Println(log)
}

// 原始sql语句查询数量
func demo5() {
	count := mapper.SelectCount4SQL("select count(*) from logs where log_type = ?", "test")
	fmt.Println(fmt.Sprintf("查询到的count: %d", count))
}

// 更新, 使用事务
func demo6() {
	qw := mapper.BuilderUpdateWrapper(&Logs{}, true)
	// Set函数可以调用多次, 但不可以一次都不调用
	qw.Set(true, "log_type", "test")
	qw.Set(true, "log_level", "INFO")
	qw.Eq(true, "id", 1)
	count := qw.Update()
	if count == 0 {
		qw.Rollback()
	} else {
		qw.Commit()
	}
	fmt.Println(fmt.Sprintf("影响行数: %d", count))
}

// 更新 - 全部 此功能慎用 会导致全表所有数据更新 , 不使用事务
func demo6_all() {
	qw := mapper.BuilderUpdateWrapper(&Logs{})
	qw.Set(true, "log_type", "test")
	qw.Set(true, "log_level", "INFO")
	count := qw.Update()
	fmt.Println(fmt.Sprintf("影响行数: %d", count))
}

// 根据ID删除数据
func demo7() {
	count := mapper.DeleteById(&Logs{}, 1)
	fmt.Println(fmt.Sprintf("影响行数: %d", count))
}

// 根据ID批量删除数据
func demo7_1() {
	var ids = []int{1, 2, 3}
	qw := mapper.BuilderUpdateWrapper(&Logs{})
	qw.DeleteByIds(ids)
}

func demo7_2() {
	var ids = []int{1, 2, 3}
	mapper.DeleteByIds(&Logs{}, ids)
}

// 根据条件删除, 使用事务
func demo8() {
	qw := mapper.BuilderUpdateWrapper(&Logs{})
	qw.Eq(true, "log_type", "test")
	count := qw.Delete()
	fmt.Println(fmt.Sprintf("影响行数: %d", count))
}

// 删除全部 此功能慎用 , 会清空表数据
func demo8_all() {
	var log Logs
	qw := mapper.BuilderUpdateWrapper(&log)
	count := qw.Delete()
	fmt.Println(fmt.Sprintf("影响行数: %d", count))
}

// 新增 默认自增主键 空值排除
func demo9() {
	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "test"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "test"

	mapper.Insert(&resLog)
	// 新增后该对象中会有新增后的id
	fmt.Println(resLog)
}

// 新增 空值排除 指定排除某些字段
// 主键是否自增根据是否给主键赋值来自动判定
func demo10() {
	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "test"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "test"
	// 排除的字段为表字段名 , 可写多个 , 写上之后 , 即使上方赋值也不会插入该值
	mapper.Insert(&resLog, "log_time", "created_date")
	// 新增后该对象中会有新增后的id
	fmt.Println(resLog)
}

func demo10_1() {

	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "test"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "test"
	wrapper := mapper.BuilderUpdateWrapper(&resLog)
	// 排除的字段为表字段名 , 可写多个 , 写上之后 , 即使上方赋值也不会插入该值
	wrapper.Insert("log_time", "created_date")
	// 新增后该对象中会有新增后的id
	fmt.Println(resLog)
}

// 新增自定义主键规则与排除字段
func demo11() {
	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "test"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "test"
	wrapper := mapper.BuilderUpdateWrapper(&resLog)
	// 第一个false为是否排除空值(false为不排除空值)
	wrapper.InsertCustom(false, "log_time", "created_date")
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

// 原生sql进行插入 - 自增主键 - 无事务
func demo13() {
	count, id := mapper.Insert4SQL(true, "insert into logs (log_time,log_type) values (?,?)", "2024-07-07 07:07:07", "system")
	fmt.Println(count, id)
}

// 有事务
func demo13_1() {
	wrapper := mapper.BuilderUpdateWrapper(nil, true)

	count, id := wrapper.Insert4SQL(true, "insert into logs (log_time,log_type) values (?,?)", "2024-07-07 07:07:07", "system")

	if count == 0 {
		wrapper.Rollback()
	} else {
		wrapper.Commit()
	}

	fmt.Println(count, id)
}

// 事务传播
func demo14() {

	var resLog Logs
	resLog.LogContent = "test"
	resLog.CreatedDate = "2024-07-07 07:07:07"
	resLog.LogLevel = "test"
	resLog.LogTime = "2024-07-07 07:07:07"
	resLog.LogType = "test"
	// 开启事务
	wrapper := mapper.BuilderUpdateWrapper(&resLog, true)
	// 排除的字段为表字段名 , 可写多个 , 写上之后 , 即使上方赋值也不会插入该值
	num := wrapper.Insert("log_time", "created_date")
	if num == 0 {
		wrapper.Rollback()
		return
	}

	qw := mapper.BuilderUpdateWrapper(&Logs{})
	// 使用事务传播
	qw.SpreadTransaction(wrapper)
	qw.Eq(true, "log_type", "test")
	count := qw.Delete()

	if count == 0 {
		qw.Rollback()
	} else {
		qw.Commit()
	}

}

func isEmpty(str string) bool {
	return len(str) == 0
}
