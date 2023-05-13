package main

import (
	"fmt"
	"geeCache/consistenthash"
	"geeCache/pb"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

const (
	defaultBasePath = "/_geecache/"
	defaultPeplicas = 50
)

// ****HTTP服务端
// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	self        string //记录自己的地址 IP:Port
	basePath    string //节点间通讯地址的前缀
	mu          sync.Mutex
	peers       *consistenthash.Map    //一致性哈希
	httpGetters map[string]*httpGetter //不同key对应不同baseURL的httpGetter
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

	//编码
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

// ******Http客户端
type httpGetter struct {
	baseURL string
}

// http的get请求,同时实现PeerGetter接口
func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	//设置url
	url := fmt.Sprintf(
		"%v%v%v",
		h.baseURL,
		//添加query参数
		url.QueryEscape(in.Group)+"/",
		url.QueryEscape(in.Key),
	)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	//客户端必须关闭body
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returnd:%v", resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//解码
	if err := proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("proto unmarshal body error:%v", err)
	}

	return nil
}

// 设置了http池中的http客户端（为每个远程节点设置）
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	//初始化哈希表
	p.peers = consistenthash.New(defaultPeplicas, nil)

	//添加哈希映射
	p.peers.Add(peers...)

	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// 实现PeerPicker接口
// 根据所需要的key值找到对应的httpGetter
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//根据一致性哈希表获取key值所对应的节点名
	peer := p.peers.Get(key)

	if peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

var _PeerGetter = (*httpGetter)(nil)
