package main

import (
	"net"

	"gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel/transport/httpupgrade"
)

type ServerConfig struct {
	ListenAddress string
}

type Transport struct {
	config *ServerConfig
}

func NewWebTunnelServerTransport(config *ServerConfig) (Transport, error) {
	return Transport{config: config}, nil
}

func (t Transport) Accept(conn net.Conn) (net.Conn, error) {
	config := &httpupgrade.Config{}
	if httpUpgradeTransport, err := httpupgrade.NewHTTPUpgradeTransport(config); err != nil {
		return nil, err
	} else {
		return httpUpgradeTransport.Server(conn)
	}
}
