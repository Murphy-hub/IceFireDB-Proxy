package proxy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/cache"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/monitor"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/bareneter"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/config"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/rediscluster"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/router"
	proxycluster "github.com/IceFireDB/IceFireDB-Proxy/pkg/router/redisCluster"
	proxynode "github.com/IceFireDB/IceFireDB-Proxy/pkg/router/redisNode"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/wconfig"
	redisclient "github.com/gomodule/redigo/redis"
)

type Proxy struct {
	Cache        *cache.Cache
	Monitor      *monitor.Monitor
	proxyCluster *rediscluster.Cluster
	proxyClient  *redisclient.Pool
	server       *bareneter.Server
	router       router.IRoutes
}

func New() (*Proxy, error) {
	p := &Proxy{}
	var err error
	if config.Get().RedisDB.Type == config.TypeNode {
		p.proxyClient = &redisclient.Pool{
			MaxIdle:     config.Get().RedisDB.ConnPoolSize,
			IdleTimeout: time.Duration(config.Get().RedisDB.ConnAliveTimeOut) * time.Second,
			Dial: func() (redisclient.Conn, error) {
				c, err := redisclient.Dial("tcp", config.Get().RedisDB.StartNodes,
					redisclient.DialConnectTimeout(time.Duration(config.Get().RedisDB.ConnTimeOut)*time.Second),
					redisclient.DialReadTimeout(time.Duration(config.Get().RedisDB.ConnReadTimeOut)*time.Second),
					redisclient.DialWriteTimeout(time.Duration(config.Get().RedisDB.ConnWriteTimeOut)*time.Second))
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redisclient.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}

		p.router = proxynode.NewRouter(p.proxyClient)
	} else {

		p.proxyCluster, err = rediscluster.NewCluster(
			&rediscluster.Options{
				StartNodes:             strings.Split(config.Get().RedisDB.StartNodes, ","),
				ConnTimeout:            time.Duration(config.Get().RedisDB.ConnTimeOut) * time.Second,
				ReadTimeout:            time.Duration(config.Get().RedisDB.ConnReadTimeOut) * time.Second,
				WriteTimeout:           time.Duration(config.Get().RedisDB.ConnWriteTimeOut) * time.Second,
				KeepAlive:              config.Get().RedisDB.ConnPoolSize,
				AliveTime:              time.Duration(config.Get().RedisDB.ConnAliveTimeOut) * time.Second,
				SlaveOperateRate:       config.Get().RedisDB.SlaveOperateRate,
				ClusterUpdateHeartbeat: config.Get().RedisDB.ClusterUpdateHeartbeat,
			})
		if err != nil {
			return nil, err
		}
		p.router = proxycluster.NewRouter(p.proxyCluster)
		wconfig.Init(context.Background(), p.proxyCluster, wconfig.Config{ReadTimeOut: time.Duration(config.Get().RedisDB.ConnReadTimeOut) * time.Second})
	}

	p.StartMonitor()

	p.router.Use(router.IgnoreCMDMiddleware(config.Get().IgnoreCMD.Enable, config.Get().IgnoreCMD.CMDList))

	p.router.Use(router.KeyMonitorMiddleware(p.Monitor, config.Get().Monitor.SlowQueryConf.SlowQueryIgnoreCMD))
	p.router.InitCMD()

	p.server = bareneter.NewServerNetwork("tcp",
		fmt.Sprintf(":%d", config.Get().Proxy.LocalPort),
		p.handle,
		p.accept,
		p.closed)
	return p, nil
}

func (p *Proxy) Run(ctx context.Context, errSignal chan error) {
	go func() {
		select {
		case <-ctx.Done():
			_ = p.server.Close()
		}
	}()
	_ = p.server.ListenServeAndSignal(errSignal)
}
