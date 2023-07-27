package main

import (
	log "github.com/sirupsen/logrus"
	"go-http-tunnel/tunnel"
)

var (
	addr = ":8080"

	remoteAddr = ":8081"
)

func main() {
	server := tunnel.NewClient(addr, remoteAddr)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
