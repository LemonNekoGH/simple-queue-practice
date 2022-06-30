package main

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisList(t *testing.T) {
	t.Run("LPUSH and LPOP", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		ctx := context.Background()

		forPush := []string{"Lemon", "Neko", "Kawaii"}
		forPop := []string{}
		// 依次将元素插入 list 列表中的左侧
		for _, elment := range forPush {
			redisClient.LPush(ctx, "list", elment)
		}
		// 测试结束时删除数据
		defer func() {
			_, err := redisClient.Del(ctx, "list").Result()
			require.NoError(err)
		}()
		// 检查 Neko 元素的位置
		index, err := redisClient.LPos(ctx, "list", "Neko", redis.LPosArgs{
			Rank: 1,
		}).Result()
		require.NoError(err)
		assert.Equal(int64(1), index)
		// 全部推出
		for _, _ = range forPush {
			popped, err := redisClient.LPop(ctx, "list").Result()
			require.NoError(err)
			forPop = append(forPop, popped)
		}
		// 检查推出之后的元素集合是否和插入之前的一致
		assert.ElementsMatch(forPush, forPop)
	})
	t.Run("LPUSH and RPOP", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		ctx := context.Background()

		forPush := []string{"Lemon", "Neko", "Kawaii"}
		forPop := []string{}
		expectForPop := []string{"Kawaii", "Neko", "Lemon"}
		// 依次将元素插入 list 列表中的左侧
		for _, elment := range forPush {
			redisClient.LPush(ctx, "list", elment)
		}
		// 测试结束时删除数据
		defer func() {
			_, err := redisClient.Del(ctx, "list").Result()
			require.NoError(err)
		}()
		// 检查 Neko 元素的位置
		index, err := redisClient.LPos(ctx, "list", "Neko", redis.LPosArgs{
			Rank: 1,
		}).Result()
		require.NoError(err)
		assert.Equal(int64(1), index)
		// 从右侧推出所有的元素
		for _, _ = range forPush {
			popped, err := redisClient.RPop(ctx, "list").Result()
			require.NoError(err)
			forPop = append(forPop, popped)
		}
		// 检查推出之后的元素集合是否和插入之前的相反
		assert.ElementsMatch(expectForPop, forPop)
	})
}
