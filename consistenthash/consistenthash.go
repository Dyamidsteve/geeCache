package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//一致性哈希

// 哈希值获取函数
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int            //虚拟节点倍数
	keys     []int          //哈希环
	hashMap  map[int]string //虚拟节点与真实节点的映射表
}

func New(replicas int, fn Hash) *Map {

	m := &Map{
		hash:     fn,
		replicas: replicas,
		keys:     make([]int, 0),
		hashMap:  make(map[int]string),
	}

	//哈希获取方法为空则使用默认方法
	if fn == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash) //添加虚拟节点到哈希环
			m.hashMap[hash] = key

		}
	}

	sort.Ints(m.keys) //升序排列
}

// 获取节点Get方法
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//计算key的哈希
	hash := int(m.hash([]byte(key)))

	//顺时针找到第一个匹配的虚拟节点下标
	fir_index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//根据下标对应的hash找到对应的真实节点
	return m.hashMap[m.keys[fir_index%len(m.keys)]]

}
