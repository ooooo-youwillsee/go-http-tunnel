package tunnel

import (
	"net"
	"time"
)

type ServerOption func(server *Server)

type ClientOption func(client *Client)

func setTCPConnOptions(conn net.Conn) {
	tcpConn := conn.(*net.TCPConn)
	//tcpConn.SetReadDeadline(time.Now().Add(30 * time.Second))
	//tcpConn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	tcpConn.SetKeepAlivePeriod(5 * time.Second)
	tcpConn.SetKeepAlive(true)
}
