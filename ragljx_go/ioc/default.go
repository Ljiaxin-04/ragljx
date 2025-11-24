package ioc

const (
	CONFIG_NAMESPACE  = "configs"
	API_NAMESPACE     = "apis"
	DEFAULT_NAMESPACE = "default"
)

var (
	// DefaultStore 默认存储
	DefaultStore = &defaultStore{
		store: []*NamespaceStore{
			newNamespaceStore(CONFIG_NAMESPACE).SetPriority(99),
			newNamespaceStore(DEFAULT_NAMESPACE).SetPriority(9),
			newNamespaceStore(API_NAMESPACE).SetPriority(-99),
		},
	}
)

// Config 用于托管配置对象的IOC空间，最先初始化
func Config() StoreUser {
	return DefaultStore.Namespace(CONFIG_NAMESPACE)
}

// Api 用于托管RestApi对象的IOC空间，最后初始化
func Api() StoreUser {
	return DefaultStore.Namespace(API_NAMESPACE)
}

// Default 默认空间，用于托管工具类，在控制器之前进行初始化
func Default() StoreUser {
	return DefaultStore.Namespace(DEFAULT_NAMESPACE)
}

