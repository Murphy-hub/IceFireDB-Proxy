package proxy

import (
	"net"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/bareneter"
	"github.com/IceFireDB/IceFireDB-Proxy/pkg/config"
	"github.com/IceFireDB/IceFireDB-Proxy/utils"
)

func (p *Proxy) accept(conn bareneter.Conn) bool {
	if config.Get().IPWhiteList.Enable {
		host, _, _ := net.SplitHostPort(conn.RemoteAddr())
		if !utils.InArray(host, config.Get().IPWhiteList.List) {
			return false
		}
	}

	p.Monitor.ConnectionGauge.Inc()
	return true
}

func (p *Proxy) closed(conn bareneter.Conn, err error) {
	p.Monitor.ConnectionGauge.Dec()
}
