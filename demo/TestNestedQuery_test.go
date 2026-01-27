package demo

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
	"testing"
)

// 定义一个测试实体
type User struct {
	Id    int    `gorm:"column:id;primaryKey;autoIncrement" tableId:"true"`
	Name  string `gorm:"column:name"`
	Age   int    `gorm:"column:age"`
	Sex   int    `gorm:"column:sex"`
	Word1 string `gorm:"column:word1"`
	Word2 string `gorm:"column:word2"`
	Key1  string `gorm:"column:key1"`
	Key2  string `gorm:"column:key2"`
	Key3  string `gorm:"column:key3"`
	Key4  string `gorm:"column:key4"`
}

// TableName 指定表名
func (u User) TableName() string {
	return "users"
}

func TestNestedQuery(t *testing.T) {
	// 构建复杂嵌套查询
	wrapper := mapper.BuilderQueryWrapper()
	result := wrapper.Eq(true, "name", "jinzhu").
		Eq(true, "sex", 1).
		//And(
		//	mapper.NewAndGroup().
		//		Eq(true, "name", "jinzhu 2").
		//		Eq(true, "age", 18),
		//).
		//And(
		//	mapper.NewOrGroup().
		//		Eq(true, "word1", "bbbb").
		//		Eq(true, "word2", "aaa"),
		//).
		Or(
			mapper.NewAndGroup().
				Eq(true, "key1", "123").
				Eq(true, "key2", 456),
		)
	//Or(
	//	mapper.NewOrGroup().
	//		Eq(true, "key3", "789").
	//		Eq(true, "key4", 666),
	//)

	// 生成SQL（这会调用queryWrapper4SQL）
	sql, values := result.GetSQLAndParams()
	fmt.Printf("Generated SQL: %s\n", sql)
	fmt.Printf("Values: %v\n", values)

	// 简单验证SQL是否生成了
	if sql != "" {
		t.Logf("Generated SQL: %s", sql)
		t.Logf("With values: %v", values)
	} else {
		t.Log("No SQL generated")
	}
}

// 测试UpdateWrapper的嵌套查询功能
func TestUpdateWrapperNestedQuery(t *testing.T) {
	// 构建复杂嵌套查询用于更新
	wrapper := &mapper.UpdateWrapper{}
	result := wrapper.Eq(true, "name", "updated_name").
		And(
			mapper.NewAndGroup().
				Eq(true, "key1", "123").
				Eq(true, "key2", 456),
		).
		Or(
			mapper.NewOrGroup().
				Eq(true, "key3", "789").
				Eq(true, "key4", 666),
		)

	// 生成SQL（这会调用queryWrapper4SQL）
	sql, values := result.GetSQLAndParams()
	fmt.Printf("Generated Update SQL: %s\n", sql)
	fmt.Printf("Update Values: %v\n", values)

	if sql != "" {
		t.Logf("Generated Update SQL: %s", sql)
		t.Logf("With update values: %v", values)
	} else {
		t.Log("No update SQL generated")
	}
}
