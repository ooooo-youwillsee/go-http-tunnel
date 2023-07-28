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
	flag.StringVar(&config, "c", "", "-c client.ini")
	flag.Parse()
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetReportCaller(true)
	ccs := tunnel.NewClientConfigsFromFile(config)
	for _, cc := range *ccs {
		go StartClient(cc.LocalAddr, cc.RemoteAddr, cc.TunnelAddr, cc.TunnelUrl)
	}
	<-tunnel.NewQuitSignal()
}

func StartClient(localAddr, proxyAddr, tunnelAddr, tunnelUrl string) {
	server := tunnel.NewClient(localAddr, proxyAddr, tunnelAddr, tunnelUrl)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
