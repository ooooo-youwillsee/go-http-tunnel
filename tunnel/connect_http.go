package tunnel

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"sync"
)

func (c *Client) connectWithHTTP(conn net.Conn) {
	// dial tunnel
	tunnelConn, err := net.Dial("tcp", c.tunnelAddr)
	if err != nil {
		log.Errorf("dial localAddr %s err %v\n", c.tunnelAddr, err)
	}
	// send request
	request, _ := http.NewRequest(http.MethodGet, c.tunnelUrl, nil)
	request.Host = c.tunnelAddr
	request.Header.Set(HEADER_REMOTE_ADDR, c.remoteAddr)
	request.Header.Set("HOST", request.Host)
	err = request.Write(tunnelConn)
	if err != nil {
		log.Error("send connect request ", err)
		return
	}
	// receive response
	response, err := http.ReadResponse(bufio.NewReader(tunnelConn), request)
	if err != nil {
		log.Error("receive connect response ", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Error("connect http tunnel err")
		return
	}

	var wg sync.WaitGroup
	// read data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConn1ToConn2(conn, tunnelConn)
	}()

	// write data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConn1ToConn2(tunnelConn, conn)
	}()
	wg.Wait()
}

func (s *Server) connectWithHTTP(w http.ResponseWriter, r *http.Request, remoteAddr string) {
	// return success
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "connected success")

	// http hijacker
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "not support hijacker")
		return
	}
	log.Infoln("http hijacker success")
	conn, _, _ := hijacker.Hijack()
	defer conn.Close()

	// dial remote addr
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Error("dial remote addr ", err)
		return
	}
	defer remoteConn.Close()

	var wg sync.WaitGroup
	// read data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConn1ToConn2(conn, remoteConn)
	}()

	// write data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConn1ToConn2(remoteConn, conn)
	}()
	wg.Wait()
}

func copyConn1ToConn2(conn1 net.Conn, conn2 net.Conn) {
	_, err := io.Copy(conn2, conn1)
	if err != nil {
		return
	}
}
