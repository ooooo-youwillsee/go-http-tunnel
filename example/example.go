package main

import "github.com/ooooo-youwillsee/go-http-tunnel/tunnel"

var (
	localAddr  = ":8080"
	remoteAddr = ":30001"
	tunnelAddr = ":8081"
)

func main() {
	server := tunnel.NewServer(tunnelAddr, "")
	go server.ListenAndServe()

	client := tunnel.NewClient(localAddr, remoteAddr, tunnelAddr, "")
	go client.ListenAndServe()

	<-tunnel.NewQuitSignal()
}
