package config

//配置

type Config struct {
	Name      string //配置文件目录
	EnvPrefix string //服务名称
}

func Init(name, envPrefix string) error {
	config := Config{
		Name:      name,
		EnvPrefix: envPrefix,
	}
	//初始化viper
	if err := config.initViper(); err != nil {
		return err
	}

	//初始化日志
	config.initZapLog()
	//初始化配置文件监控
	config.watchViperConfig()
	return nil
}

