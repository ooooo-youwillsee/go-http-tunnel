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
	flag.StringVar(&config, "c", "", "-c server.ini")
	flag.Parse()
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetReportCaller(true)
	scs := tunnel.NewServerConfigsFrom(config)
	for _, sc := range *scs {
		go startServer(sc.Addr, sc.Url)
	}
	<-tunnel.NewQuitSignal()
}

func startServer(addr string, url string) {
	server := tunnel.NewServer(addr, url)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
