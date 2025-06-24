package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container 是 IOC 容器的核心结构体
type Container struct {
	mu         sync.RWMutex                  // 读写锁，保护并发安全
	instances  map[string]interface{}        // 实例缓存表：name → 实例（实现单例）
	providers  map[string]func() interface{} // 构造函数注册表：name → 构造函数
	lifecycles []Lifecycle                   // 实现了生命周期接口的实例列表（统一 Start/Stop）
}

// New 创建并初始化一个新的 IOC 容器实例
func New() *Container {
	return &Container{
		instances:  make(map[string]interface{}),        // 初始化实例缓存
		providers:  make(map[string]func() interface{}), // 初始化构造器映射
		lifecycles: make([]Lifecycle, 0),                // 初始化生命周期管理列表
	}
}

// Provide 注册一个构造函数，供容器用于延迟实例化
// 参数 name 是唯一标识，constructor 是无参构造函数（如 func() *Logger）
func (c *Container) Provide(name string, constructor interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 使用反射检查 constructor 是否是函数类型
	val := reflect.ValueOf(constructor)
	if val.Kind() != reflect.Func {
		panic("constructor must be a function")
	}

	// 构造函数必须是：无参、返回一个值的形式，如 func() *T
	if val.Type().NumIn() != 0 || val.Type().NumOut() != 1 {
		panic("constructor must be of type func() T")
	}

	// 将构造函数封装为闭包，注册到 providers 中
	c.providers[name] = func() interface{} {
		return val.Call(nil)[0].Interface() // 调用函数并取第一个返回值
	}
}

// Get 获取服务实例，并将其实例化结果注入到传入的目标变量（target）中
// target 必须是一个指针，如：&logger
func (c *Container) Get(target interface{}) {
	ptrVal := reflect.ValueOf(target)

	// 必须是指针类型才能注入（因为需要修改原始变量的值）
	if ptrVal.Kind() != reflect.Ptr {
		panic("target must be a pointer")
	}

	// 获取指针所指向变量的类型名作为 key（如 *Logger → "main.Logger"）
	name := ptrVal.Type().Elem().String()

	// 尝试从实例缓存中读取（读锁）
	c.mu.RLock()
	instance, ok := c.instances[name]
	c.mu.RUnlock()

	if ok {
		// 如果已缓存，直接注入（*target = instance）
		ptrVal.Elem().Set(reflect.ValueOf(instance))
		return
	}

	// 没有缓存，则从构造函数表中查找
	c.mu.RLock()
	provider, ok := c.providers[name]
	c.mu.RUnlock()

	if !ok {
		// 如果连构造函数都没注册，报错
		panic(fmt.Sprintf("no provider found for %s", name))
	}

	// 调用构造函数，构造新的实例（懒加载）
	instance = provider()

	// 写入缓存（写锁）
	c.mu.Lock()
	c.instances[name] = instance
	c.mu.Unlock()

	// 注入实例到目标变量中（*target = instance）
	ptrVal.Elem().Set(reflect.ValueOf(instance))

	// 如果实例实现了 Lifecycle 接口，加入生命周期管理
	if l, ok := instance.(Lifecycle); ok {
		c.mu.Lock()
		c.lifecycles = append(c.lifecycles, l)
		c.mu.Unlock()
	}
}

// InitAll 遍历所有实现了 Lifecycle 的实例，依次调用其 Start() 方法
func (c *Container) InitAll() {
	for _, l := range c.lifecycles {
		_ = l.Start() // 忽略错误，可按需修改为记录或中止
	}
}

// StopAll 遍历所有实现了 Lifecycle 的实例，依次调用其 Stop() 方法
func (c *Container) StopAll() {
	for _, l := range c.lifecycles {
		_ = l.Stop() // 忽略错误
	}
}
