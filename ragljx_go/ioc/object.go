package ioc

import (
	"context"
	"fmt"
)

// ObjectImpl 对象基础实现，其他对象可以继承此结构体
type ObjectImpl struct {
}

func (i *ObjectImpl) Init() error {
	return nil
}

func (i *ObjectImpl) Name() string {
	return ""
}

func (i *ObjectImpl) Close(ctx context.Context) {
}

func (i *ObjectImpl) Priority() int {
	return 0
}

// ObjectWrapper 对象包装器
type ObjectWrapper struct {
	Name     string
	Priority int
	Value    Object
}

func NewObjectWrapper(obj Object) *ObjectWrapper {
	name := obj.Name()
	if name == "" {
		name = fmt.Sprintf("%T", obj)
	}
	return &ObjectWrapper{
		Name:     name,
		Priority: obj.Priority(),
		Value:    obj,
	}
}

