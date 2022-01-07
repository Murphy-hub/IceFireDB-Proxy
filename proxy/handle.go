package proxy

import (
	"errors"
	"fmt"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/RedSHandle"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/bareneter"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/codis/credis"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/router"
	"github.com/sirupsen/logrus"
)

func (p *Proxy) handle(conn bareneter.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	localConn := conn.NetConn()
	localWriteHandle := RedSHandle.NewWriterHandle(localConn)
	decoder := credis.NewDecoderSize(localConn, 1024)
	for {
		resp, err := decoder.Decode()
		if err != nil {
			/*if err.Error() != io.EOF.Error() && strings.Index(err.Error(), net.ErrClosed.Error()) == -1 {
				logrus.Errorf("RESP协议解码失败:%v", err)
			}*/
			return
		}

		if resp.Type != credis.TypeArray {
			_ = router.WriteError(localWriteHandle, fmt.Errorf(router.ErrUnknownCommand, "cmd"))
			return
		}

		respCount := len(resp.Array)

		if respCount < 1 {
			_ = router.WriteError(localWriteHandle, fmt.Errorf(router.ErrArguments, "cmd"))
			return
		}

		if resp.Array[0].Type != credis.TypeBulkBytes {
			_ = router.WriteError(localWriteHandle, router.ErrCmdTypeWrong)
			return
		}

		commandArgs := make([]interface{}, respCount)
		for i := 0; i < respCount; i++ {
			commandArgs[i] = resp.Array[i].Value
		}

		err = p.router.Handle(localWriteHandle, commandArgs)

		if err != nil {
			if errors.Is(err, router.ErrLocalWriter) || errors.Is(err, router.ErrLocalFlush) {
				return
			}
			_ = router.WriteError(localWriteHandle, err)
			logrus.Errorf("redis命令执行错误:%s , %v", commandArgs, err)
			return
		}
	}
}
