/*
 *
 *  * Licensed to the Apache Software Foundation (ASF) under one or more
 *  * contributor license agreements.  See the NOTICE file distributed with
 *  * this work for additional information regarding copyright ownership.
 *  * The ASF licenses this file to You under the Apache License, Version 2.0
 *  * (the "License"); you may not use this file except in compliance with
 *  * the License.  You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

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
