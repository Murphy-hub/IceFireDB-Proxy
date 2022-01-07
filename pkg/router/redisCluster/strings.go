package redisCluster

import (
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/rediscluster"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/router"
)

func (r *Router) cmdMGET(s *router.Context) error {
	batch := r.redisCluster.NewBatch(10)
	argLen := len(s.Args)
	for i := 1; i < argLen; i++ {
		err := batch.Put(r.redisCluster.IsSlaveOperate(), "GET", s.Args[i])
		if err != nil {
			return err
		}
	}
	var err error
	reply, err := rediscluster.Values(r.redisCluster.RunBatch(batch))
	if err != nil && err != rediscluster.ErrNil {
		return router.WriteError(s.Writer, err)
	}

	if len(reply) == 1 {
		if err, ok := reply[0].(error); ok {
			return router.WriteError(s.Writer, err)
		}
	}
	s.Reply = reply
	return router.RecursivelyWriteObjects(s.Writer, reply...)
}

func (r *Router) cmdDEL(s *router.Context) error {
	argLen := len(s.Args)

	if argLen == 2 {
		reply, err := rediscluster.Int64(r.redisCluster.Do(s.Cmd, s.Args[1]))
		if err != nil {
			return router.WriteError(s.Writer, err)
		}
		return router.WriteInt(s.Writer, reply)
	}

	batch := r.redisCluster.NewBatch(10)
	for i := 1; i < argLen; i++ {
		err := batch.Put(false, "DEL", s.Args[i])
		if err != nil {
			return err
		}
	}
	reply, err := rediscluster.Values(r.redisCluster.RunBatch(batch))
	if err != nil && err != rediscluster.ErrNil {
		return router.WriteError(s.Writer, err)
	}

	var delCount int64
	for _, valInterface := range reply {
		if count, ok := valInterface.(int64); ok {
			delCount += count
		}
	}
	s.Reply = delCount
	return router.WriteInt(s.Writer, delCount)
}

func (r *Router) cmdEXISTS(s *router.Context) error {
	argLen := len(s.Args)

	if argLen == 2 {
		reply, err := rediscluster.Int64(r.redisCluster.Do(s.Cmd, s.Args[1]))
		if err != nil {
			return router.WriteError(s.Writer, err)
		}
		return router.WriteInt(s.Writer, reply)
	}

	batch := r.redisCluster.NewBatch(10)
	for i := 1; i < argLen; i++ {
		err := batch.Put(false, s.Cmd, s.Args[i])
		if err != nil {
			return err
		}
	}
	reply, err := rediscluster.Values(r.redisCluster.RunBatch(batch))
	if err != nil && err != rediscluster.ErrNil {
		return router.WriteError(s.Writer, err)
	}

	var existsCount int64
	for _, valInterface := range reply {
		if count, ok := valInterface.(int64); ok {
			existsCount += count
		}
	}
	s.Reply = existsCount
	return router.WriteInt(s.Writer, existsCount)
}
