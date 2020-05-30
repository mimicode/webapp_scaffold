package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

func (c *Config) initViper() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./conf")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(c.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "-"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

//监控配置文件改动
func (c *Config) watchViperConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("config is changed", zap.String("file", in.String()))
	})
}
