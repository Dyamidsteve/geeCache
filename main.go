package main

import (
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

type server int

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("hellow world"))
}

func main() {
	NewGroup("scores", GetterFunc(func(key string) ([]byte, error) {
		log.Println("[slowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), 2<<10)

}

// var set = make(map[int]bool, 0)
// var lck sync.Mutex
// func printOnce(num int) {
// 	lck.Lock()
// 	if _, exist := set[num]; !exist {
// 		fmt.Println(num)
// 	}
// 	set[num] = true
// 	lck.Unlock()
// }

// func main() {
// 	for i := 0; i < 5; i++ {
// 		go printOnce(100)
// 	}
// 	time.Sleep(time.Second)
// }
