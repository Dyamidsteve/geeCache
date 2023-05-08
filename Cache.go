package main

import (
	"geeCache/lru"
	"sync"
)

// 实例化lru，封装add和get方法，并添加了锁处理并发
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) Add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//延迟初始化(Lazy Initialization)，
	//一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。
	//主要用于提高性能，并减少程序内存要求。
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Add(key, &val)

}

func (c *cache) Get(key string) (bv ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if val, ok := c.lru.Get(key); ok {

		return val.(ByteView), true
	}

	return
}
