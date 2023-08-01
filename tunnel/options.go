package tunnel

import (
	"net"
	"time"
)

type ServerOption func(server *Server)

type ClientOption func(client *Client)

func ServerWithToken(token string) ServerOption {
	return func(server *Server) {
		server.token = token
	}
}

func ClientWithToken(token string) ClientOption {
	return func(client *Client) {
		client.token = token
	}
}

func ClientWithMode(mode ConnectMode) ClientOption {
	return func(client *Client) {
		client.mode = mode
	}
}

func setTCPConnOptions(conn net.Conn) {
	tcpConn := conn.(*net.TCPConn)
	//tcpConn.SetReadDeadline(time.Now().Add(30 * time.Second))
	//tcpConn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	tcpConn.SetKeepAlivePeriod(5 * time.Second)
	tcpConn.SetKeepAlive(true)
}
