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
	tunnelAddrs := tcpProxies.GetTunnelAddrs()
	for _, addr := range tunnelAddrs {
		go startServer(addr)
	}
	<-tunnel.NewQuitSignal()
}

func startServer(addr string) {
	server := tunnel.NewServer(addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
