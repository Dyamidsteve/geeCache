package main

import (
	"errors"
	"fmt"
	"geeCache/pb"
	"geeCache/singleflight"
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
	name      string     //组名
	getter    Getter     //本地data获取器
	mainCache cache      //缓存
	peers     PeerPicker //远程节点获取器(接口化，不限制http、tcp、udp)
	// use singleflight.Group to make sure that
	//each key is only fetched once
	loader *singleflight.Group
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
		peers:     nil,
		loader:    &singleflight.Group{},
	}

	groups[nam] = group

	return group

}

// 从全局变量groups中获取*Group
func GetGroup(name string) *Group {
	//只使用读锁R，与RW锁互斥
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// 注册远端节点获取器
func (g *Group) RegisterPeers(pk PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = pk
}

// 找到对应的节点，并请求获取数据
func (g *Group) load(key string) (val ByteView, err error) {
	//调用loader.Do，确保并发场景下相同的key，load过程只会调用一次
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			//查看是否对应远程节点
			if peer, ok := g.peers.PickPeer(key); ok {
				if val, err = g.getFromPeer(key, peer); err == nil {
					return val, nil
				}
				log.Println("[geeCache] failed to get from peer", err)
			}
		}

		//远程节点没有,只能找本地节点
		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}

	return

}

// 从远程节点获取数据
func (g *Group) getFromPeer(key string, peer PeerGetter) (ByteView, error) {
	msg_request := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	msg_resp := &pb.Response{}
	err := peer.Get(msg_request, msg_resp)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: msg_resp.Value}, nil
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
