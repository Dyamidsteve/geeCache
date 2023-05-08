package main

//Getter接口获取data
type Getter interface {
	Get(key string) ([]byte, error)
}

//****这种接口型函数使用场景更好，既能够接收匿名函数，也能够接收普通函数，还能接收实现该方法的结构体
//定义方法类型实现Getter接口
type GetterFunc func(key string) ([]byte, error)

//实现Getter的Get方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
