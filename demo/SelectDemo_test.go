package demo

import (
	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
	"testing"
)

func TestSelectCount(t *testing.T) {

	registerDemo()

	count := mapper.SelectCount[user](mapper.BuilderQueryWrapper())

	logger.Info("查询所有用户数量为 count: {}", count)

}

func TestSelectCount2(t *testing.T) {

	registerDemo()

	count := mapper.SelectCount[user](mapper.BuilderQueryWrapper().
		Ge(true, "age", 20))

	logger.Info("查询年龄大于20岁的数量为 count: {}", count)

}

func TestSelectId(t *testing.T) {
	registerDemo()
	users := mapper.SelectById[user](1)
	logger.Info("查询id为1的 user: {}", users)
}

func TestSelectList(t *testing.T) {
	registerDemo()
	users := mapper.SelectList[user](mapper.BuilderQueryWrapper())
	logger.Info("查询所有的 user: {}", users)
}

func TestSelectList2(t *testing.T) {
	registerDemo()
	users := mapper.SelectList[user](mapper.BuilderQueryWrapper().
		Ge(true, "age", 20))
	logger.Info("查询年龄大于20岁的 user: {}", users)
}

func TestSelectList3(t *testing.T) {
	registerDemo()
	ids := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	users := mapper.SelectList[user](mapper.BuilderQueryWrapper().
		InInt32(true, "id", ids),
	)
	logger.Info("查询指定ids的 user: {}", users)
}
