package demo

import (
	"github.com/PurpleScorpion/go-sweet-orm/v3/logger"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
	"testing"
)

func TestUpdate(t *testing.T) {

	registerDemo()

	count := mapper.Update[user](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", 1).
		Set(true, "age", 51),
	)

	logger.Info("更新成功 count: {}", count)

}
