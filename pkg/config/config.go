package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/rediscluster"
	"strings"
)

var (
	ErrConfigNotInit       = errors.New("Config not init")
	ErrDuplicateInitConfig = errors.New("Duplicate init config！")
)

const (
	TypeNode    = "node"
	TypeCluster = "cluster"
)

// 代理的数据链路中不要直接使用config配置，防止race
// 全局配置
var (
	_config          *Config
	_ruleRedisClient *rediscluster.Cluster
)

func InitConfig() error {
	if _config != nil {
		return ErrDuplicateInitConfig
	}

	err := viper.Unmarshal(&_config)
	if err != nil {
		return err
	}
	if _config.IgnoreCMD.Enable && len(_config.IgnoreCMD.CMDList) > 0 {
		CmdToUpper(_config.IgnoreCMD.CMDList)
	}

	if _config.Monitor.SlowQueryConf.Enable && len(_config.Monitor.SlowQueryConf.SlowQueryIgnoreCMD) > 0 {
		CmdToUpper(_config.Monitor.SlowQueryConf.SlowQueryIgnoreCMD)
	}
	return nil
}


func Get() *Config {
	return _config
}

func CmdToUpper(list []string) {
	for k := range list {
		list[k] = strings.ToUpper(list[k])
	}
}