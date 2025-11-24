package ioc

import "context"

// Store IOC容器存储接口
type Store interface {
	StoreUser
	StoreManage
}

// StoreUser 用户使用的存储接口
type StoreUser interface {
	// Registry 对象注册
	Registry(obj Object)
	// Get 对象获取
	Get(name string) Object
	// List 打印对象列表
	List() []string
	// Len 数量统计
	Len() int
}

// StoreManage 存储管理接口
type StoreManage interface {
	// LoadFromEnv 从环境变量中加载对象配置
	LoadFromEnv(prefix string) error
	// LoadFromFile 从配置文件加载配置
	LoadFromFile(path string) error
}

// Object IOC对象接口，需要注册到IOC空间托管的对象需要实现的方法
type Object interface {
	// Init 对象初始化，初始化对象的属性
	Init() error
	// Name 对象的名称，根据名称可以从空间中取出对象
	Name() string
	// Priority 对象优先级，根据优先级控制对象初始化的顺序，数值越大优先级越高
	Priority() int
	// Close 对象的销毁，服务关闭时调用
	Close(ctx context.Context)
}

