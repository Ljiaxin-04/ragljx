package ioc

var (
	isLoaded bool
)

// LoadConfigRequest 配置加载请求
type LoadConfigRequest struct {
	// 环境变量配置
	ConfigEnv *ConfigEnv
	// 文件配置方式
	ConfigFile *ConfigFile
}

type ConfigFile struct {
	Enabled bool
	Path    string
}

type ConfigEnv struct {
	Enabled bool
	Prefix  string
}

func NewLoadConfigRequest() *LoadConfigRequest {
	return &LoadConfigRequest{
		ConfigEnv: &ConfigEnv{
			Enabled: true,
		},
		ConfigFile: &ConfigFile{
			Enabled: true,
			Path:    "config/application.yaml",
		},
	}
}

// ConfigIocObject 配置IOC对象
func ConfigIocObject(req *LoadConfigRequest) error {
	if isLoaded {
		return nil
	}

	// 优先加载环境变量
	if req.ConfigEnv.Enabled {
		err := DefaultStore.LoadFromEnv(req.ConfigEnv.Prefix)
		if err != nil {
			return err
		}
	}

	// 再加载配置文件
	if req.ConfigFile.Enabled {
		err := DefaultStore.LoadFromFile(req.ConfigFile.Path)
		if err != nil {
			return err
		}
	}

	// 初始化对象
	err := DefaultStore.InitIocObject()
	if err != nil {
		return err
	}

	isLoaded = true
	return nil
}

