package main

import (
	"flag"
	"github.com/ooooo-youwillsee/go-http-tunnel/tunnel"
	log "github.com/sirupsen/logrus"
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
		go StartClient(cc)
	}
	<-tunnel.NewQuitSignal()
}

func StartClient(cc *tunnel.ClientConfig) {
	server := tunnel.NewClient(cc)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
