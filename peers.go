package main

import "geeCache/pb"

//节点选择器
type PeerPicker interface {
	//通过key值返回对应节点的getter来获取数据
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//节点数据获取器接口，这里使用protobuf装载数据
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
