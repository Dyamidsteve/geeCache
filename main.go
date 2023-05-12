package main

import (
	"fmt"
	"log"
	"net/http"
)

// type server int

// func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	log.Println(r.URL.Path)
// 	w.Write([]byte("hellow world"))
// }

func main() {
	db := map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
	//添加group,并设置数据源Getter方法
	NewGroup("scores", GetterFunc(func(key string) ([]byte, error) {
		log.Println("[slowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), 2<<10)

	addr := "localhost:9999"
	peer := NewHTTPPool(addr)

	log.Printf("Http Server start on %s", addr)
	log.Fatal(http.ListenAndServe(addr, peer))

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
