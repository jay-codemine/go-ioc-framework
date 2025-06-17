package container

import (
	"reflect"
	"sync"
)

// ContextKey 是一个内部私有类型，避免与用户冲突
type ContextKey struct{}

// ScopedContainer 用于每个请求上下文的临时容器
type ScopedContainer struct {
	parent     *Container
	scopedData sync.Map // 用于保存每个请求作用域下的实例
}

// NewScopedContainer 创建一个请求级容器
func NewScopedContainer(parent *Container) *ScopedContainer {
	return &ScopedContainer{
		parent: parent,
	}
}

// Set 设置请求上下文中的一个值
func (s *ScopedContainer) Set(name string, val interface{}) {
	s.scopedData.Store(name, val)
}

// Get 获取请求上下文中的值，如果没有则回退到全局容器
func (s *ScopedContainer) Get(target interface{}) {
	name := typeNameOfPtr(target)

	// scoped 优先
	if val, ok := s.scopedData.Load(name); ok {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(val))
		return
	}

	// 回退到全局容器
	s.parent.Get(target)
}

// typeNameOfPtr 用于识别类型名称（作为 key）
func typeNameOfPtr(target interface{}) string {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		panic("target must be a pointer")
	}
	return t.Elem().String()
}
