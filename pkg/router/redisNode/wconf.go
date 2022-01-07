/**
 * @Author: zhangchao
 * @Description:
 * @Date: 2021/6/16 4:06 下午
 */
package redisNode

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/RedSHandle"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/rediscluster"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/router"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/wconfig"
	"github.com/sirupsen/logrus"
)

const (
	cmdTimeout = "BLOCK"
	cmdConfigs = "CONFIGS"
	cmdPrefix  = "PREFIX"
)

func (r *Router) cmdWCONFIG(s *router.Context) error {
	return fmt.Errorf(router.ErrUnknownCommand, "wconfig")
}

func (r *Router) cmdWCONFIGNotReady(s *router.Context) error {
	args := s.Args
	if len(args) < 5 {
		return errors.New("ERR wrong number of arguments for 'wconfig' command")
	}

	streamKey := string(args[1].([]byte))
	consumerID := string(args[2].([]byte))
	var (
		configIDs   []string
		updateAt    []int64
		prefix      string
		timeoutInMS int = 20
	)
	args = args[3:]
paramLoop:
	for len(args) > 0 {
		switch strings.ToUpper(string(args[0].([]byte))) {
		case cmdTimeout:
			timeout := string(args[1].([]byte))
			i, err := strconv.Atoi(string(timeout))
			if err != nil {
				logrus.Warnf("watch config block param parse err: %s", err.Error())
				return errors.New("ERR wrong block param for 'wconfig' command")
			}
			timeoutInMS = i
			args = args[2:]
		case cmdPrefix:
			prefix = string(args[1].([]byte))
			args = args[2:]
		case cmdConfigs:
			args = args[1:]
			if len(args)%2 == 1 {
				return errors.New("ERR wrong number of CONFIGS for 'wconfig' command")
			}
			count := len(args) / 2
			configIDs = make([]string, count)
			updateAt = make([]int64, count)
			for i := 0; i < count; i++ {
				configIDs[i] = prefix + string(args[i].([]byte))
				updateAtStr := string(args[count+i].([]byte))
				if updateAtStr == "$" {
					continue
				}
				updateAt[i], _ = strconv.ParseInt(updateAtStr, 10, 64)
			}
			break paramLoop
		default:
			if len(args) > 1 {
				args = args[2:]
			} else {
				args = nil
			}
		}
	}
	if len(configIDs) == 0 {
		return errors.New("ERR CONFIGS not defined for 'wconfig' command")
	}

	for i := 0; i < len(updateAt); i++ {
		if updateAt[i] > 0 {
			changed := r.checkConfigChanged(configIDs[i], updateAt[i])
			if changed {
				c, err := r.getConfig(configIDs[i])
				if err != nil {
					c = nil
				}
				return router.RecursivelyWriteObjects(s.Writer, configIDs[i], c)
			}
		}
	}

	w := wconfig.AddConfigWatcher(streamKey, consumerID, configIDs...)
	if w == nil {
		return router.WriteInt(s.Writer, 0)
	}

	time.AfterFunc(time.Millisecond*time.Duration(timeoutInMS), func() {
		_ = w.Stop()
	})

	value, err := w.Next()
	if err != nil || len(value) == 0 {
		return WriteConfig(s.Writer, nil)
	}
	w.Stop()

	return WriteConfig(s.Writer, value)
}

func (r *Router) checkConfigChanged(configKey string, updateAt int64) bool {
	timestamp, err := rediscluster.Int64(r.Do("HGET", configKey, "timestamp"))
	if err != nil {
		return false
	}
	return timestamp > updateAt
}

func (r *Router) getConfig(configKey string) ([]interface{}, error) {
	return rediscluster.Values(r.Do("HGETALL", configKey))
}

func WriteConfig(local *RedSHandle.WriterHandle, reply map[string]*wconfig.ConfigResponse) (err error) {
	if len(reply) == 0 {
		err = local.RecursivelyWriteObjects()
	} else {
		var c *wconfig.ConfigResponse
		for _, con := range reply {
			c = con
		}
		err = router.RecursivelyWriteObjects(local, c.ConfigID, c.ConfigInfo)
	}
	if err != nil {
		logrus.Error("客户端写入数据错误", err)
		return router.ErrLocalWriter
	}
	err = local.Flush()
	if err != nil {
		logrus.Error("客户端flush错误", err)
		return router.ErrLocalFlush
	}
	return nil
}
