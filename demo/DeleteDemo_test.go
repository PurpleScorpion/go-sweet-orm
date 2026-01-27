package demo

import (
	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
	"testing"
)

func TestDelete(t *testing.T) {
	registerDemo()
	count := mapper.Delete[user](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", 1),
	)
	logger.Info("删除成功: count: {}", count)
}

func TestDeleteById(t *testing.T) {
	registerDemo()
	count := mapper.DeleteById[user](2, nil)
	logger.Info("删除成功: count: {}", count)
}

func TestDeleteByIds(t *testing.T) {
	registerDemo()
	ids := []int{3, 4, 5}
	count := mapper.DeleteByIds[user](ids, nil)
	logger.Info("删除成功: count: {}", count)
}
