package redisNode

import (
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/router"
	"github.com/gomodule/redigo/redis"
)

var (
	pongReply = "PONG"
	okReply   = "OK"
)

func (r *Router) cmdCOMMAND(s *router.Context) error {
	return router.WriteObjects(s.Writer, nil)
}

func (r *Router) cmdPING(s *router.Context) error {
	s.Reply = pongReply
	return router.WriteSimpleString(s.Writer, pongReply)
}

func (r *Router) cmdCMDEXEC(s *router.Context) error {
	var err error
	s.Reply, err = r.Do(s.Cmd, s.Args[1:]...)
	if err != nil && err != redis.ErrNil {
		_ = router.WriteError(s.Writer, err)
		return nil
	}

	if s.Reply == nil {
		return router.WriteBulk(s.Writer, nil)
	}

	switch val := s.Reply.(type) {
	case error:
		return router.WriteError(s.Writer, val)
	case int64:
		return router.WriteInt(s.Writer, val)
	case []byte:
		return router.WriteBulk(s.Writer, val)
	case string:
		return router.WriteSimpleString(s.Writer, val)
	case []interface{}:
		if len(val) == 1 {
			if err, ok := val[0].(error); ok {
				return router.WriteError(s.Writer, err)
			}
		}
		return router.RecursivelyWriteObjects(s.Writer, val...)
	default:
		return router.WriteObjects(s.Writer, s.Reply)
	}
}

/*func (r *Router) cmdCMDEXEC(s *router.Context) error {


	var rt redis.ReplyType
	var err error
	rt, s.Reply, err = r.client.Do(r.ctx, s.Args...).FullResult()
	if rt == redis.ReplyTypeError {
		s.Reply = err
		err = nil
	}

	if err != nil && err != redis.Nil {

		return err
	}

	if s.Reply == nil {
		return router.WriteBulk(s.Writer, nil)
	}

	switch val := s.Reply.(type) {
	case error:
		return router.WriteError(s.Writer, val)
	case int64:
		return router.WriteInt(s.Writer, val)
	case []byte:
		return router.WriteBulk(s.Writer, val)
	case string:
		if rt == redis.ReplyTypeString {
			return router.WriteBulk(s.Writer, []byte(val))
		}
		return router.WriteSimpleString(s.Writer, val)
	case []interface{}:
		if len(val) == 1 {
			if err, ok := val[0].(error); ok {
				return router.WriteError(s.Writer, err)
			}
		}
		return router.RecursivelyWriteObjects(s.Writer, val...)
	default:
		return router.WriteObjects(s.Writer, s.Reply)
	}
}*/

func (r *Router) cmdQUIT(s *router.Context) error {
	s.Reply = okReply
	return router.WriteSimpleString(s.Writer, okReply)
}
