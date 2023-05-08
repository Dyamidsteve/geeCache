package lru

import "container/list"

// LRU 最近最少使用，被访问的移到队尾，淘汰时淘汰队首
// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64 //最大使用内存
	nBytes   int64 //已使用内存

	ll    *list.List               //双向链表
	cache map[string]*list.Element //用key指向链表元素的表

	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数，可以为 nil。
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

func NewCache(Mbytes int64, onEvictedFunc func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  Mbytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvictedFunc,
	}
}

// 获取缓存，并将该缓存放到队尾
func (c *Cache) Get(key string) (val Value, ok bool) {
	//1、找到该key对应的element
	if element, ok := c.cache[key]; ok {
		//2、将该element移到队尾
		c.ll.MoveToFront(element)
		ent := element.Value.(*entry)
		return ent.value, true

	}
	return
}

//淘汰最老的缓存
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()

	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//更新cache这个map
		delete(c.cache, kv.key)

		//缓存计数减少
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			//回调函数不为空则触发
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

//添加或修改数据
func (c *Cache) Add(key string, val Value) {
	//最大内存小于等于0则无法添加
	if c.maxBytes <= 0 {
		return
	}

	if ele, ok := c.cache[key]; ok {
		//已经存在
		c.ll.MoveToFront(ele)
		ent := ele.Value.(*entry)

		//修改当前使用的内存
		c.nBytes += int64(val.Len() - ent.value.Len())
		ele.Value = val //赋值

	} else {
		//将其添加到队尾
		entry := c.ll.PushFront(&entry{key: key, value: val})
		//为cache中map添加该映射
		c.cache[key] = entry
		//增加内存消耗
		c.nBytes += int64(len(key)) + int64(val.Len())
	}

	//添加或修改后若内存不够淘汰最久未被访问的
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
