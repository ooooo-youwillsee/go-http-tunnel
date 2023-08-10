package main

import (
	"github.com/ooooo-youwillsee/go-http-tunnel/tunnel"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	localAddr  = ":8111"
	remoteAddr = ":30001"
	tunnelAddr = ":8112"
)

func main() {
	server := tunnel.NewServer(tunnelAddr, "")
	go server.ListenAndServe()

	client := tunnel.NewClient(localAddr, remoteAddr, tunnelAddr, "", tunnel.ClientWithSMux("true"))
	go client.ListenAndServe()

	go testServer(remoteAddr)

	<-tunnel.NewQuitSignal()
}

func testServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/hello", hello)
	log.Infoln("listen addr ", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalln(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "HOST = "+r.Header.Get("HOST"))
}

func index(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "index")
}
