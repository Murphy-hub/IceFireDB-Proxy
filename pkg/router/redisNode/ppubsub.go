package redisNode

import (
	"errors"
	"fmt"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/ppubsub"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/router"
)

func (r *Router) cmdPpub(s *router.Context) error {
	args := s.Args
	if len(args) != 3 {
		return errors.New("ERR wrong number of arguments for 'ppub' command")
	}

	topicName := string(args[1].([]byte))
	message := string(args[2].([]byte))
	err := ppubsub.Pub(topicName, message)
	if err != nil {
		return errors.New("ERR pub:" + err.Error())
	}
	return router.WriteInt(s.Writer, 1)
}

func (r *Router) cmdPsub(s *router.Context) error {
	args := s.Args
	if len(args) != 2 {
		return errors.New("ERR wrong number of arguments for 'psub' command")
	}
	topicName := string(args[1].([]byte))
	ps, err := ppubsub.Sub(topicName)
	if err != nil {
		return errors.New("ERR sub:" + err.Error())
	}
	for {
		msg := <-ps.Inbound
		//var data []interface{}
		//err := json.Unmarshal([]byte(msg.Message), &data)
		//if err != nil {
		//	logrus.Errorf("subscribe error: %v", err)
		//	continue
		//}
		fmt.Println(msg)
		router.WriteBulkStrings(s.Writer, []string{msg.Message})
		//data = make([]interface{}, 2)
		//data[0] = msg.Message
		//err = r.Sync(data)
		//if err != nil {
		//	logrus.Errorf("subscribe sync error: %v", err)
		//	continue
		//}
	}
	//return
}
