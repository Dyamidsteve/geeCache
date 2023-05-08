package main

import (
	"reflect"
	"testing"
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
