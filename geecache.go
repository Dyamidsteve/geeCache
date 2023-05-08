package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

// Getter接口获取data
type Getter interface {
	Get(key string) ([]byte, error)
}

// ****这种接口型函数使用场景更好，既能够接收匿名函数，也能够接收普通函数，还能接收实现该方法的结构体
// 定义方法类型实现Getter接口
type GetterFunc func(key string) ([]byte, error)

// 实现Getter的Get方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(nam string, get Getter, cacheBytes int64) *Group {
	if get == nil {
		panic("nil getter")
	}
	//使用RW锁
	mu.Lock()
	defer mu.Unlock()

	group := &Group{
		name:      nam,
		getter:    get,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[nam] = group

	return group

}

func GetGroup(name string) *Group {
	//只使用读锁R，与RW锁互斥
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// 从缓存和数据源中获取数据，数据源也找不到则返回非空的error
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}
	//先从缓存中寻找
	if val, ok := g.mainCache.Get(key); ok {
		log.Println("[GeeCache] hit")
		return val, nil
	}

	//再从外存(数据源)中寻找
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		fmt.Println("group getter Get error:", err)
		return ByteView{}, err
	}

	bv := ByteView{b: bytes}

	//将从数据源获取的数据存到cache中
	g.mainCache.Add(key, bv)

	return bv, nil

}
