// 我们通常搭配 WATCH命令来执行事务操作。从使用WATCH命令监视某个 key 开始，直到执行EXEC命令的这段时间里，
// 如果有其他用户抢先对被监视的 key 进行了替换、更新、删除等操作，那么当用户尝试执行EXEC的时候，事务将失败并返回一个错误，
// 用户可以根据这个错误选择重试事务或者放弃事务。

// Watch方法接收一个函数和一个或多个key作为参数。
Watch(fn func(*Tx) error, keys ...string) error

// 下面的代码片段演示了 Watch 方法搭配 TxPipelined 的使用示例。
// watchDemo 在key值不变的情况下将其值+1
func watchDemo(ctx context.Context, key string) error {
	return rdb.Watch(ctx, func(tx *redis.Tx) error {
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		// 假设操作耗时5秒
		// 5秒内我们通过其他的客户端修改key，当前事务就会失败
		time.Sleep(5 * time.Second)
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n+1, time.Hour)
			return nil
		})
		return err
	}, key)
}

// 将上面的函数执行并打印其返回值，如果我们在程序运行后的5秒内修改了被 watch 的 key 的值，那么该事务操作失败，返回redis: transaction failed错误。


// 最后我们来看一个 go-redis 官方文档中使用 GET 、SET和WATCH命令实现一个 INCR 命令的完整示例。

const routineCount = 100 // 此处rdb为初始化的redis连接客户端

// 设置5秒超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// increment 是一个自定义对key进行递增（+1）的函数
// 使用 GET + SET + WATCH 实现，类似 INCR
increment := func(key string) error {
	txf := func(tx *redis.Tx) error {
		// 获得当前值或零值
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// 实际操作（乐观锁定中的本地操作）
		n++

		// 仅在监视的Key保持不变的情况下运行
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// pipe 处理错误情况
			pipe.Set(ctx, key, n, 0)
			return nil
		})
		return err
	}

	// 最多重试100次
	for retries := routineCount; retries > 0; retries-- {
		err := rdb.Watch(ctx, txf, key)
		if err != redis.TxFailedErr {
			return err
		}
		// 乐观锁丢失
	}
	return errors.New("increment reached maximum number of retries")
}

// 开启100个goroutine并发调用increment
// 相当于对key执行100次递增
var wg sync.WaitGroup
wg.Add(routineCount)
for i := 0; i < routineCount; i++ {
	go func() {
		defer wg.Done()

		if err := increment("counter3"); err != nil {
			fmt.Println("increment error:", err)
		}
	}()
}
wg.Wait()

n, err := rdb.Get(ctx, "counter3").Int()
fmt.Println("最终结果：", n, err)



