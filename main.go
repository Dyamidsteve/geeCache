package main

import (
	"fmt"
	"sync"
	"time"
)

var set = make(map[int]bool, 0)
var lck sync.Mutex

func printOnce(num int) {
	lck.Lock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
	lck.Unlock()
}

func main() {
	for i := 0; i < 5; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}
