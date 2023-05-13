package singleflight

import "sync"

//处理缓存击穿

//call代表正在进行中，或已经结束的请求。
type call struct {
	wg  sync.WaitGroup //sync.WaitGroup 锁避免重入。
	val interface{}
	err error
}

//Group 是 singleflight 的主数据结构，
//管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

/*
Do 方法，接收 2 个参数，第一个参数是 key，第二个参数是一个函数 fn。
Do 的作用就是，针对相同的 key，
无论 Do 被调用多少次，函数 fn 都只会被调用一次，
等待 fn 调用结束了，返回返回值或错误。
*/
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//若已经有对应的请求
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		//等待所有任务完成
		c.wg.Wait()

		//将该请求的数据，错误返回
		return c.val, c.err
	}
	//生成一个空值call请求指针
	c := new(call)

	//开始新的任务
	c.wg.Add(1)

	//为group中的callMap添加新的请求
	g.m[key] = c

	//这里先添加请求再释放锁可以保证其他groutine
	//进入有对应请求的程序块内
	g.mu.Unlock()

	//最先获得锁的才能执行fn，保证同时只能执行一次fn
	c.val, c.err = fn()

	//给请求赋值后再结束自己的任务，告知其他groutine不用等待可以返回了
	c.wg.Done()

	//请求任务结束后，删除该请求，腾出内存，不需要一直保存
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
