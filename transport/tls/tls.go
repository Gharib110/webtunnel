package tls

import (
	"crypto/tls"
	"errors"
	"net"
)

type Config struct {
	ServerName string
}

func NewTLSTransport(config *Config) (Transport, error) {
	return Transport{kind: "tls", serverName: config.ServerName}, nil
}

type Transport struct {
	kind       string
	serverName string
}

func (t Transport) Client(conn net.Conn) (net.Conn, error) {
	switch t.kind {
	case "tls":
		conf := &tls.Config{ServerName: t.serverName}
		return tls.Client(conn, conf), nil
	}
	return nil, errors.New("unknown kind")
}
