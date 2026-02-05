package demo

import (
	"testing"

	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
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

func TestSelectPage(t *testing.T) {
	registerDemo()
	//pageUtils := mapper.BuilderPageUtils(1, 2, mapper.BuilderQueryWrapper())
	//page := mapper.Page[user](pageUtils)
	page := mapper.Page[user](
		mapper.BuilderPageUtils(1, 2,
			mapper.BuilderQueryWrapper(),
		),
	)

	logger.Info("当前页: {}", page.Current)
	logger.Info("总页数: {}", page.TotalPage)
	logger.Info("总记录数: {}", page.TotalCount)

	list := page.List
	for _, u := range list {
		logger.Info("姓名: {} , 年龄: {}", u.UserName, u.Age)
	}

}

type UserVO struct {
	Name   string `json:"Name"`
	AgeAAA int    `json:"AgeAAA"`
}

func TestSelectPageVO(t *testing.T) {
	registerDemo()
	//pageUtils := mapper.BuilderPageUtils(1, 2, mapper.BuilderQueryWrapper())
	//page := mapper.Page[user](pageUtils)
	page := mapper.Page[user](
		mapper.BuilderPageUtils(1, 2,
			mapper.BuilderQueryWrapper(),
		),
	)

	var userVos []UserVO

	list := page.List
	for _, u := range list {
		userVos = append(userVos, UserVO{
			Name:   u.UserName,
			AgeAAA: u.Age,
		})
	}

	pageVO := mapper.ConvertPageData[user, UserVO](page, userVos)

	logger.Info("转换结果: {}", pageVO)
}
