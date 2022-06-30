package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

const queueName = "MyQueue"

var wg = sync.WaitGroup{}

// initRedisClient 初始化 Redis 客户端
func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	fmt.Println("Redis 客户端初始化完毕")
}

// 队列
func queue() {
	fmt.Println("队列开始运行")
	stop := false
	for !stop {
		// 获取任务，如果没有获取到会阻塞
		task, err := redisClient.BRPop(ctx, 0, queueName).Result()
		if err != nil {
			fmt.Printf("获取任务时失败，将重试 %s\n", err.Error())
		} else {
			if task[1] == "stop" {
				fmt.Println("收到停止指令")
				stop = true
			} else {
				fmt.Printf("正在执行 %s 任务\n", task[1])
				// 延迟 5 秒模拟耗时任务
				time.Sleep(time.Second * 5)
				fmt.Printf("%s 任务执行完成\n", task[1])
			}
		}
	}
	fmt.Println("队列结束运行")
	// 告诉协程组，运行结束
	wg.Done()
}

// 添加任务到队列中去
func addTask(content string) {
	fmt.Printf("正在将 %s 添加到队列中\n", content)
	length, err := redisClient.LPush(ctx, queueName, content).Result()
	if err != nil {
		fmt.Printf("将 %s 添加到队列时出错：%s\n", content, err.Error())
		return
	}
	fmt.Printf("%s 已被添加到队列中，目前队列中有 %d 个任务\n", content, length)

}

func main() {
	initRedisClient()
	// 创建协程组
	wg.Add(2)
	// 启动队列
	go queue()
	// 添加任务
	go func() {
		for i := 0; i <= 10; i++ {
			addTask(fmt.Sprintf("TASK %d", i))
			// 间隔两秒添加任务
			// time.Sleep(2 * time.Second)
		}
		// 停止队列
		addTask("stop")
		wg.Done()
	}()
	// 等待协程结束
	wg.Wait()
}
