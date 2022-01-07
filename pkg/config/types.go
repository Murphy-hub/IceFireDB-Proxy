package config

import "github.com/IceFireDB/IceFireDB-Proxy/pkg/monitor"

type Config struct {
	Proxy                  ProxyS                `mapstructure:"proxy"`
	RedisDB                RedisDBS              `mapstructure:"redisdb"`
	PprofDebug             PprofDebugS           `mapstructure:"pprof_debug"`
	Log                    LogS                  `mapstructure:"log"`
	Discovery              DiscoveryS            `mapstructure:"discovery"`
	IPWhiteList            IPWhiteListS          `mapstructure:"ip_white_list"`
	Cache                  CacheS                `mapstructure:"cache"`
	Monitor                MonitorS              `mapstructure:"monitor"`
	PrometheusExporterConf *monitor.ExporterConf `mapstructure:"prometheus_exporter"`
	IgnoreCMD              IgnoreCMDS            `mapstructure:"ignore_cmd"`
}

type ProxyS struct {
	// 配置最后更新时间
	LocalPort  int  `mapstructure:"local_port" json:"local_port"`   // 代理时本地监听的端口
	EnableMTLS bool `mapstructure:"enable_mtls" json:"enable_mtls"` // 集群节点，多个，分割
}

type PprofDebugS struct {
	Enable bool   `mapstructure:"enable"`
	Port   uint16 `mapstructure:"port"`
}

type LogS struct {
	Level     string `mapstructure:"level"`
	cmd       bool   `mapstructure:"cmd"`
	OutPut    string `mapstructure:"output"`
	Key       string `mapstructure:"key"`
	RedisHost string `mapstructure:"redis_host"`
	RedisPort int    `mapstructure:"redis_port"`
}

type DiscoveryS struct {
	Address   string `mapstructure:"address"`   // discovery remote address
	Heartbeat int    `mapstructure:"heartbeat"` // 上报到msp心跳周期
	Env       string `mapstructure:"env"`       // discovery Env
	Region    string `mapstructure:"region"`
	Zone      string `mapstructure:"zone"`
	NodeIP    string `mapstructure:"node_ip"`
	PodIP     string `mapstructure:"pod_ip"`
	Hostname  string `mapstructure:"hostname"`
	Endpoint  string `mapstructure:"endpoint"`
}

// ServiceInfoS is core Data struct of Miscro Service configuration
type ServiceInfoS struct {
	Env         string            `mapstructure:"env"`
	ServiceName string            `mapstructure:"service_name"`
	UniqueID    string            `mapstructure:"unique_id"`
	MemLimits   string            `mapstructure:"mem_limits"`
	CPULimits   string            `mapstructure:"cpu_limits"`
	Version     string            `mapstructure:"version"`
	MetaData    map[string]string `mapstructure:"meta_data"` // sidecar增加metadata数据要同步增加msp显示
}

// cache配置
type CacheS struct {
	Enable            bool `mapstructure:"enable"`
	MaxItemsSize      int  `mapstructure:"max_items_size"`     // cache 中最大的items数量
	DefaultExpiration int  `mapstructure:"default_expiration"` // cache kv默认的过期时间（单位：毫秒）
	CleanupInterval   int  `mapstructure:"cleanup_interval"`   // cache 过期kv清理周期（单位：秒）
}

type HotKeyConfS struct {
	Enable                  bool `mapstructure:"enable" json:"enable"`                                       // 是否开启hotkey监控
	MonitorJobInterval      int  `mapstructure:"monitor_job_interval" json:"monitor_job_interval"`           // 监控任务时间轮策略：单次休眠时间 （单位：秒 ）
	MonitorJobLifetime      int  `mapstructure:"monitor_job_lifetime" json:"monitor_job_lifetime"`           // 监控任务时间轮策略：单次监控任务维持时间 （单位：秒）
	SecondHotThreshold      int  `mapstructure:"second_hot_threshold" json:"second_hot_threshold"`           // hotkey的鉴定阀值，在监控之后，如果发现key的每秒处理量高于这个阀值，则判定为hotkey
	SecondIncreaseThreshold int  `mapstructure:"second_increase_threshold" json:"second_increase_threshold"` // 命令处理量增长速率（单位： 次/秒） 如果中间件发现命令处理量增长速率高于这个阀值，则立刻开启hotkey监控
	LruSize                 int  `mapstructure:"lru_size" json:"-"`                                          // hotkey最大存储空间，存储按照lru策略
	EnableCache             bool `mapstructure:"enable_cache" json:"enable_cache"`                           // 是否开启中间件缓存
	MaxCacheLifeTime        int  `mapstructure:"max_cache_life_time" json:"max_cache_life_time"`             // 缓存条目最大存活时间
}

type BigKeyConfS struct {
	Enable           bool `mapstructure:"enable" json:"enable"`                           // 是否开启bigkey监控
	KeyMaxBytes      int  `mapstructure:"key_max_bytes" json:"key_max_bytes"`             // redis key最大长度限制（单位：byte）
	ValueMaxBytes    int  `mapstructure:"value_max_bytes" json:"value_max_bytes"`         // redis value最大长度限制（单位：byte）
	LruSize          int  `mapstructure:"lru_size" json:"-"`                              // bigkey最大存储空间，存储按照lru策略
	EnableCache      bool `mapstructure:"enable_cache" json:"enable_cache"`               // 是否开启中间件缓存
	MaxCacheLifeTime int  `mapstructure:"max_cache_life_time" json:"max_cache_life_time"` // 缓存条目最大存活时间
}

type SlowQueryConfS struct {
	Enable             bool     `mapstructure:"enable" json:"enable"`
	SlowQueryThreshold int      `mapstructure:"slow_query_threshold" json:"slow_query_threshold"`   // slow query的时间阀值：高于此阀值则进入slow query lru内存
	MaxListSize        int      `mapstructure:"max_list_size" json:"-"`                             // hotkey最大存储空间，存储按照队列策略
	SlowQueryIgnoreCMD []string `mapstructure:"slow_query_ignore_cmd" json:"slow_query_ignore_cmd"` // 慢查忽略的key
}

type IgnoreCMDS struct {
	Enable  bool     `mapstructure:"enable" json:"enable"`
	CMDList []string `mapstructure:"cmd_list" json:"cmd_list"` // 慢查忽略的key
}

type MonitorS struct {
	HotKeyConf    HotKeyConfS    `mapstructure:"hotkey"`
	BigKeyConf    BigKeyConfS    `mapstructure:"bigkey"`
	SlowQueryConf SlowQueryConfS `mapstructure:"slowquery"`
}

// RedisClusterConf is redis cluster configure options
type RedisDBS struct {
	Type                   string `mapstructure:"type"` // node、cluster
	StartNodes             string `mapstructure:"start_nodes"`
	ConnTimeOut            int    `mapstructure:"conn_timeout"`             // Connection timeout parameter of cluster nodes Unit: ms
	ConnReadTimeOut        int    `mapstructure:"conn_read_timeout"`   // Cluster node read timeout parameter Unit: ms
	ConnWriteTimeOut       int    `mapstructure:"conn_write_timeout"` // Cluster node write timeout parameter Unit: ms
	ConnAliveTimeOut       int    `mapstructure:"conn_alive_timeout"` // Cluster node TCP idle survival time Unit: seconds
	ConnPoolSize           int    `mapstructure:"conn_pool_size"`         // The size of the TCP connection pool for each node in the cluster
	SlaveOperateRate       int    `mapstructure:"slave_operate_rate"`
	ClusterUpdateHeartbeat int    `mapstructure:"cluster_update_heartbeat"`
}

type IPWhiteListS struct {
	Enable bool     `json:"enable"`
	List   []string `json:"list"`
}
