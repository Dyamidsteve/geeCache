package lru

import (
	"fmt"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestLruGet(t *testing.T) {
	lru := NewCache(1024, func(s string, v Value) {
		fmt.Println("test")
	})

	lru.Add("a", String("aaa"))
	if ele, ok := lru.Get("a"); !ok || string(ele.(String)) != "aaa" {
		t.Fatalf("cache hit key='a' failed;\n")
	}

	if _, ok := lru.Get("b"); ok {
		t.Fatalf("cache miss key='b' failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	lru := NewCache(int64(len(k1+k2+k3+v1+v2+v3)), nil)

	lru.Add(k1, String(v1)) //4byte
	lru.Add(k2, String(v2)) //4byte
	lru.Add(k3, String(v3)) //4byte

	len1 := lru.Len()
	//fmt.Println(lru.Get("ddd"))
	lru.Add("dd", String("bb")) //4byte

	//当添加{"dd","bb"}这个数据后，最老的{"k1","v1"}会被删除
	if _, ok := lru.Get("k1"); ok || lru.Len() != len1 {
		t.Fatalf("cache not automically removeOldest due to the exceed of maxBytes")
	}
}

// 测试回调函数
func TestOnEvicted(t *testing.T) {
	a := 1
	lru := NewCache(100, func(s string, v Value) {
		//fmt.Println("key:", s, "val:", v)
		a = 100
	})

	lru.Add("a", String("b"))
	lru.RemoveOldest()

	if a != 100 {
		t.Fatalf("回调函数未被调用")
	}

}
