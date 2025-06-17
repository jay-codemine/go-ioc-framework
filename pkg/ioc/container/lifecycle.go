package container

// Lifecycle 定义统一生命周期接口
type Lifecycle interface {
	Start() error
	Stop() error
}
