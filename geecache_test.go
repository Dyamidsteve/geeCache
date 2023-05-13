package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var (
	Db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
)

func TestGetter(t *testing.T) {
	var test1 Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")

	if v, _ := test1.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("回调函数失败")
	}

}

func TestGroupGet(t *testing.T) {
	//记录访问每个key的次数
	loadCounts := make(map[string]int, len(Db))
	//这里没有直接用grouo中的cache而是自定义map来获取数据
	g := NewGroup("g1", GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDb] search key:", key)
			if val, ok := Db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}

				// 访问次数++
				loadCounts[key]++
				return []byte(val), nil
			}

			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)

	//bv, err := g.Get("Tom")
	for key, val := range Db {
		if v, err := g.Get(key); err != nil || v.String() != val {
			t.Fatalf("not find [key:%s,value:%s]", key, val)
		}
		if _, err := g.Get(key); err != nil || loadCounts[key] > 1 {
			t.Fatalf("cache %s miss", key)
		} // cache hit
	}

	if v, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty,but get %s", v.String())
	}
}
