package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go-http-tunnel/tunnel"
)

var (
	config string
)

func init() {
	flag.StringVar(&config, "c", "", "-c config.ini")
	flag.Parse()
}

func main() {
	tcpProxies := tunnel.NewTcpProxiesFromFile(config)
	for _, tcpProxy := range *tcpProxies {
		go StartClient(tcpProxy.LocalAddr, tcpProxy.RemoteAddr, tcpProxy.TunnelAddr)
	}
	<-tunnel.NewQuitSignal()
}

func StartClient(localAddr, proxyAddr, tunnelAddr string) {
	server := tunnel.NewClient(localAddr, proxyAddr, tunnelAddr)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
