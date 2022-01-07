/**
 * @Author: zhangchao
 * @Description:
 * @Date: 2021/6/16 4:37 下午
 */
package wconfig

import (
	"container/list"
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/rediscluster"
)

const (
	configID  = "configID"
	updatedAt = "updatedAt"
	key       = "key"
	encoding  = "encoding"
)

var defaultConfig = Config{
	prefetchCount: 10,
	ReadTimeOut:   time.Minute,
}

var defaultStreamStore *streamStore

func Init(ctx context.Context, client *rediscluster.Cluster, c Config) {
	defaultStreamStore = NewStreamStore(ctx, client, c)
}

func AddConfigWatcher(stream, consumerID string, configIDs ...string) *Watcher {
	return defaultStreamStore.AddConfigWatcher(stream, consumerID, configIDs...)
}

type Config struct {
	prefetchCount int
	ReadTimeOut   time.Duration
}

type streamStore struct {
	sync.RWMutex
	ctx             context.Context
	client          *rediscluster.Cluster
	isClosing       bool
	configs         map[string]*list.List
	connectingLocks map[string]*sync.Mutex
	streams         map[string]struct{}
	c               Config
}

func NewStreamStore(ctx context.Context, client *rediscluster.Cluster, c Config) *streamStore {
	s := &streamStore{
		ctx:             ctx,
		client:          client,
		configs:         map[string]*list.List{},
		connectingLocks: map[string]*sync.Mutex{},
		streams:         map[string]struct{}{},
		c:               c,
	}
	var empty Config
	if c == empty {
		s.c = defaultConfig
	}
	return s
}

func (r *streamStore) AddConfigWatcher(stream, consumerID string, configIDs ...string) *Watcher {
	if r.ctx.Err() != nil {
		return nil
	}

	w := &Watcher{
		exit:    make(chan interface{}),
		updates: make(chan map[string]*ConfigResponse, 1),
	}
	consumers := map[*list.List]*list.Element{}
	r.Lock()
	for i := 0; i < len(configIDs); i++ {
		watchers := r.configs[configIDs[i]]
		if watchers == nil {
			watchers = list.New()
			r.configs[configIDs[i]] = watchers
		}

		e := watchers.PushBack(w)
		consumers[watchers] = e
	}
	r.Unlock()
	go func() {
		<-w.exit
		r.Lock()
		for watchers, e := range consumers {
			watchers.Remove(e)
		}
		r.Unlock()
	}()

	r.ensureStreamOpen(stream)
	return w
}

func (r *streamStore) ensureStreamOpen(stream string) {
	if r.steamExist(stream) {
		return
	}
	r.Lock()
	destinationLock, hasConnected := r.connectingLocks[stream]
	if !hasConnected {
		destinationLock = &sync.Mutex{}
		r.connectingLocks[stream] = destinationLock
	}
	r.Unlock()
	destinationLock.Lock()
	r.RLock()
	_, exist := r.streams[stream]
	if exist {
		r.RUnlock()
		destinationLock.Unlock()
		return
	}
	r.RUnlock()

	go r.pollNewMessagesLoop(stream)
	r.Lock()
	r.streams[stream] = struct{}{}
	destinationLock.Unlock()
	delete(r.connectingLocks, stream)
	r.Unlock()
	logrus.Infof("建立stream监听：%s", stream)
}

func (r *streamStore) pollNewMessagesLoop(stream string) {
	id := "$"
	for {

		if r.ctx.Err() != nil {
			return
		}

		streams, err := r.client.Do("XREAD", "COUNT", 1, "BLOCK", strconv.Itoa(int(r.c.ReadTimeOut/time.Millisecond)), "STREAMS", stream, id)
		if err != nil {
			if !errors.Is(err, rediscluster.ErrNil) {
				err = errors.Unwrap(err)
				netErr, ok := err.(*net.OpError)
				if ok && netErr.Timeout() {
					continue
				}
				logrus.Warnf("redis streams: error reading from stream %s: %s", stream, err)
				time.Sleep(200 * time.Millisecond)
			}
			continue
		}
		var data map[string]string
		if s, ok := streams.([]interface{}); ok && len(s) > 0 {
			id, data = parseXreadStream(s)
		}
		if c := data[configID]; c != "" {
			logrus.Infof("感应到配置变化configID ：%s", c)
			res := make(map[string]*ConfigResponse, 1)
			config, err := r.getConfig(c)
			if err != nil {
				config = nil
			}
			res[c] = &ConfigResponse{ConfigID: c, ConfigInfo: config}
			watchers := r.copyWatchers(c)
			for e := watchers.Front(); e != nil; e = e.Next() {
				w := e.Value.(*Watcher)
				select {
				case w.updates <- res:
				default:
				}
			}
		}
	}
}

func (r *streamStore) getConfig(configKey string) ([]interface{}, error) {
	return rediscluster.Values(r.client.Do("HGETALL", configKey))
}

func (r *streamStore) copyWatchers(key string) *list.List {
	watcherList := list.New()
	r.RLock()
	watchers := r.configs[key]
	if watchers != nil {
		watcherList.PushBackList(watchers)
	}
	r.RUnlock()

	return watcherList
}

func (r *streamStore) steamExist(stream string) bool {
	r.RLock()
	defer r.RUnlock()
	if _, exist := r.streams[stream]; exist {
		return true
	}
	return false
}

func (r *streamStore) DelSteam(stream string) {
	r.RLock()
	defer r.RUnlock()
	delete(r.streams, stream)
}

func parseXreadStream(s []interface{}) (id string, data map[string]string) {
	if len(s) > 0 {
		item := s[0].([]interface{})[1].([]interface{})[0].([]interface{})
		id = string(item[0].([]byte))
		d := item[1].([]interface{})
		if len(d)%2 != 0 || len(d) == 0 {
			return
		}
		data = make(map[string]string)
		for i := 0; i < len(d); i += 2 {
			key, _ := d[i].([]byte)
			value, _ := d[i+1].([]byte)
			data[string(key)] = string(value)
		}
	}
	return
}
