package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container 是 IOC 容器的核心结构
type Container struct {
	mu         sync.RWMutex
	instances  map[string]interface{}        // 实例缓存（单例）
	providers  map[string]func() interface{} // 构造函数注册表
	lifecycles []Lifecycle                   // 生命周期对象列表
}

// Lifecycle 定义可被统一启动/销毁的模块
type Lifecycle interface {
	Start() error
	Stop() error
}

// New 创建一个新的容器实例
func New() *Container {
	return &Container{
		instances:  make(map[string]interface{}),
		providers:  make(map[string]func() interface{}),
		lifecycles: make([]Lifecycle, 0),
	}
}

// Provide 注册一个服务构造函数
func (c *Container) Provide(name string, constructor interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 校验构造函数是否可调用
	val := reflect.ValueOf(constructor)
	if val.Kind() != reflect.Func {
		panic("constructor must be a function")
	}

	// 支持无参构造函数（如 func() *MyService）
	if val.Type().NumIn() != 0 || val.Type().NumOut() != 1 {
		panic("constructor must be of type func() T")
	}

	c.providers[name] = func() interface{} {
		return val.Call(nil)[0].Interface()
	}
}

// Get 获取服务实例（使用指针类型传入目标变量）
func (c *Container) Get(target interface{}) {
	ptrVal := reflect.ValueOf(target)
	if ptrVal.Kind() != reflect.Ptr {
		panic("target must be a pointer")
	}

	name := ptrVal.Type().Elem().String()

	c.mu.RLock()
	instance, ok := c.instances[name]
	c.mu.RUnlock()

	if ok {
		ptrVal.Elem().Set(reflect.ValueOf(instance))
		return
	}

	// 没有缓存，调用构造函数
	c.mu.RLock()
	provider, ok := c.providers[name]
	c.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("no provider found for %s", name))
	}

	instance = provider()

	// 缓存实例
	c.mu.Lock()
	c.instances[name] = instance
	c.mu.Unlock()

	ptrVal.Elem().Set(reflect.ValueOf(instance))

	// 如果实现了 Lifecycle 接口，加入容器管理
	if l, ok := instance.(Lifecycle); ok {
		c.mu.Lock()
		c.lifecycles = append(c.lifecycles, l)
		c.mu.Unlock()
	}
}

// InitAll 启动所有生命周期模块（调用 Start）
func (c *Container) InitAll() {
	for _, l := range c.lifecycles {
		_ = l.Start()
	}
}

// StopAll 停止所有生命周期模块（调用 Stop）
func (c *Container) StopAll() {
	for _, l := range c.lifecycles {
		_ = l.Stop()
	}
}
