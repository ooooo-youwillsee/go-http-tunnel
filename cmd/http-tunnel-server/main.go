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
	flag.StringVar(&config, "c", "", "-c server.ini")
	flag.Parse()
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetReportCaller(true)
	scs := tunnel.NewServerConfigsFrom(config)
	for _, sc := range *scs {
		go startServer(sc)
	}
	<-tunnel.NewQuitSignal()
}

func startServer(sc *tunnel.ServerConfig) {
	server := tunnel.NewServer(sc)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
