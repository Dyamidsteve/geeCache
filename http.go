package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	defaultBasePath = "/_geecache"
)

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	self     string //记录自己的地址 IP:Port
	basePath string //节点间通讯地址的前缀
}

func NewHTTPPool(sf string) *HTTPPool {
	return &HTTPPool{
		self:     sf,
		basePath: defaultBasePath,
	}
}

// Log函数
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 实现Handler接口
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}

	//记录请求方式和请求路径
	p.Log("%s %s", r.Method, r.URL.Path)

	// /basepath/groupname/key
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	//找到对应group
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
