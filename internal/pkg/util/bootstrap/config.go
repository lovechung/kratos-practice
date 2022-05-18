package bootstrap

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return file.NewSource(filePath)
}

// NewConfigProvider 创建一个配置
func NewConfigProvider(configPath string) config.Config {
	return config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
		),
	)
}
