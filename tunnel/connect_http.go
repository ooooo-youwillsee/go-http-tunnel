package tunnel

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
)

func (c *Client) connectWithHTTP() net.Conn {
	// dial tunnel
	tunnelConn, err := net.Dial("tcp", c.config.TunnelAddr)
	if err != nil {
		log.Errorf("dial TunnelAddr %s, err: %v", c.config.TunnelAddr, err)
		return nil
	}
	// send request
	request, _ := http.NewRequest(http.MethodConnect, c.config.TunnelUrl, nil)
	request.Host = c.config.TunnelAddr
	request.Header.Set("HOST", request.Host)
	c.setHeader(&request.Header)
	err = request.Write(tunnelConn)
	if err != nil {
		log.Error("send connect request ", err)
		return nil
	}
	// receive response
	response, err := http.ReadResponse(bufio.NewReader(tunnelConn), request)
	if err != nil {
		log.Error("receive connect response ", err)
		return nil
	}
	if response.StatusCode != http.StatusOK {
		log.Error("connect http tunnel err")
		return nil
	}
	return tunnelConn
}

func (s *Server) connectWithHTTP(w http.ResponseWriter, r *http.Request) net.Conn {
	if r.Method != http.MethodGet {
		log.Errorf("auth method '%s' is not supported", r.Method)
		return nil
	}
	// return success
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "connected success")

	// http hijacker
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "not support hijacker")
		return nil
	}
	log.Infoln("http hijacker success")
	conn, _, _ := hijacker.Hijack()
	return conn
}
