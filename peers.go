package main

//节点选择器
type PeerPicker interface {
	//通过key值返回对应节点的getter来获取数据
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//节点数据获取器接口
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
