package demo

import (
	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
	"testing"
)

func TestInsert1(t *testing.T) {
	registerDemo()
	var u user
	u.UserName = "test"
	u.Age = 18

	count := mapper.Insert[user](&u, nil)
	logger.Info("添加成功: count: {} , user: {}", count, u)
}

// 排除age字段
func TestInsert2(t *testing.T) {
	registerDemo()
	var u user
	u.UserName = "test2"
	u.Age = 18
	count := mapper.Insert[user](&u,
		mapper.BuilderUpdateWrapper(false).
			SetExcludeField("age"))
	logger.Info("添加成功: count: {} , user: {}", count, u)
}

// 排除空值
func TestInsert3(t *testing.T) {
	registerDemo()
	var u user
	u.UserName = ""
	u.Age = 18
	count := mapper.Insert[user](&u,
		mapper.BuilderUpdateWrapper(false).
			SetExcludeEmpty(true))
	logger.Info("添加成功: count: {} , user: {}", count, u)
}
