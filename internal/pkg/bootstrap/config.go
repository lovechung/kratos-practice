package bootstrap

import (
	consulConfig "github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	consulApi "github.com/hashicorp/consul/api"
	"strings"
)

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return file.NewSource(filePath)
}

// NewConsulConfigSource 创建一个远程配置源 - Consul
func NewConsulConfigSource(host, token, key string) config.Source {
	consulClient, err := consulApi.NewClient(&consulApi.Config{
		Address: host,
		Token:   token,
	})
	if err != nil {
		panic(err)
	}

	consulSource, err := consulConfig.New(consulClient,
		consulConfig.WithPath(getConfigKey(key, true)),
	)
	if err != nil {
		panic(err)
	}

	return consulSource
}

// getConfigKey 获取合法的配置名
func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}

// NewConfigProvider 创建一个配置
func NewConfigProvider(configPath string) config.Config {
	return config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
		),
	)
}
