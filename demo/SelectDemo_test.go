package demo

import (
	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
	"testing"
)

func TestSelectCount(t *testing.T) {

	registerDemo()

	count := mapper.SelectCount[user](mapper.BuilderQueryWrapper())

	logger.Info("查询所有用户数量为: %d", count)

}

func TestSelectCount2(t *testing.T) {

	registerDemo()

	count := mapper.SelectCount[user](mapper.BuilderQueryWrapper().
		Ge(true, "age", 20))

	logger.Info("查询年龄大于20岁的数量为: %d", count)

}

func TestSelectId(t *testing.T) {
	registerDemo()
	users := mapper.SelectById[user](1)
	logger.Info("查询id为1的user: %s", toString(users))
}

func TestSelectList(t *testing.T) {
	registerDemo()
	users := mapper.SelectList[user](mapper.BuilderQueryWrapper())
	logger.Info("查询所有的user: %s", toString(users))
}

func TestSelectList2(t *testing.T) {
	registerDemo()
	users := mapper.SelectList[user](mapper.BuilderQueryWrapper().
		Ge(true, "age", 20))
	logger.Info("查询年龄大于20岁的user: %s", toString(users))
}
