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
	tunnelUrl  = "/"
)

func main() {
	sc := &tunnel.ServerConfig{
		TunnelAddr: tunnelAddr,
		TunnelUrl:  tunnelUrl,
		Token:      "",
	}
	server := tunnel.NewServer(sc)
	go server.ListenAndServe()

	cc := &tunnel.ClientConfig{
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
		TunnelAddr: tunnelAddr,
		TunnelUrl:  tunnelUrl,
		Token:      "",
		IsSmux:     true,
		Mode:       tunnel.CONNECT_WEBSOCKET,
	}
	client := tunnel.NewClient(cc)
	go client.ListenAndServe()

	go testServer(remoteAddr)

	<-tunnel.NewQuitSignal()
}

func testServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/hello", hello)
	log.Infoln("listen test server addr ", addr)
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
