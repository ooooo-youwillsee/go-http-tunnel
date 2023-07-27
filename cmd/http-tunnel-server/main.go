package main

import (
	log "github.com/sirupsen/logrus"
	"go-http-tunnel/tunnel"
)

var (
	addr = ":8081"

	remoteAddr = ":30001"
)

func main() {
	server := tunnel.NewServer(addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
