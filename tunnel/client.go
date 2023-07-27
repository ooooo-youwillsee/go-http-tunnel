package tunnel

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type Client struct {
	localAddr  string
	remoteAddr string
	tunnelAddr string
	tcpProxies *TCPProxies
}

func NewClient(localAddr, remoteAddr, tunnelAddr string, options ...ClientOption) *Client {
	c := &Client{
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		tunnelAddr: tunnelAddr,
	}
	if localAddr == "" || remoteAddr == "" || tunnelAddr == "" {
		panic("localAddr or remoteAddr or tunnelAddr is empty")
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *Client) ListenAndServe() error {
	l, err := net.Listen("tcp", c.localAddr)
	if err != nil {
		log.Fatalln("listen localAddr err", err)
		return err
	}
	defer l.Close()
	log.Printf("listen localAddr %s\n", c.localAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error("accept err", err)
			continue
		}
		go c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn net.Conn) {
	defer conn.Close()
	// dial tunnel
	tunnelConn, err := net.Dial("tcp", c.tunnelAddr)
	if err != nil {
		log.Errorf("dial localAddr %s err %v\n", c.tunnelAddr, err)
		return
	}
	// connect
	success := c.Connect(tunnelConn, c.remoteAddr)
	if !success {
		return
	}
	// copy data
	errCh := CopyConn(conn, tunnelConn)
	<-errCh
}

func (c *Client) Connect(tunnelConn net.Conn, remoteAddr string) bool {
	// send request
	request, _ := http.NewRequest(http.MethodConnect, URL_CONNECT, nil)
	request.Header.Set(HEADER_REMOTE_ADDR, remoteAddr)
	err := request.Write(tunnelConn)
	if err != nil {
		log.Error("send connect request err", err)
		return false
	}
	// receive response
	response, err := http.ReadResponse(bufio.NewReader(tunnelConn), request)
	if err != nil {
		log.Error("receive connect response err", err)
		return false
	}
	if response.StatusCode != http.StatusOK {
		log.Error("connect http tunnel err")
		return false
	}
	return true
}
