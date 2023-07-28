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
	tunnelUrl  string
}

func NewClient(localAddr, remoteAddr, tunnelAddr, tunnelUrl string, options ...ClientOption) *Client {
	if localAddr == "" || remoteAddr == "" || tunnelAddr == "" {
		panic("localAddr or remoteAddr or tunnelAddr is empty")
	}
	if tunnelUrl == "" {
		tunnelUrl = URL_CONNECT
	}
	c := &Client{
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		tunnelAddr: tunnelAddr,
		tunnelUrl:  tunnelUrl,
	}
	for _, option := range options {
		option(c)
	}
	log.Infof("NewClient localAddr[%s], remoteAddr[%s], tunnelAddr[%s], tunnelUrl[%s]", localAddr, remoteAddr, tunnelAddr, tunnelUrl)
	return c
}

func (c *Client) ListenAndServe() error {
	l, err := net.Listen("tcp", c.localAddr)
	if err != nil {
		log.Fatalln("listen localAddr err", err)
		return err
	}
	defer l.Close()
	log.Infof("listen localAddr %s", c.localAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error("accept err", err)
			continue
		}
		cc := newClientConn(conn)
		go c.handleConn(cc)
	}
}

func (c *Client) handleConn(conn *clientConn) {
	defer conn.Close()
	// connect
	tunnelConn := c.Connect()
	if tunnelConn == nil {
		return
	}
	// copy data
	errCh := CopyConn(conn, tunnelConn)
	<-errCh
}

func (c *Client) Connect() net.Conn {
	// dial tunnel
	tunnelConn, err := net.Dial("tcp", c.tunnelAddr)
	if err != nil {
		log.Errorf("dial localAddr %s err %v\n", c.tunnelAddr, err)
		return nil
	}
	// send request
	host, _ := splitAddr(c.tunnelAddr)
	request, _ := http.NewRequest(http.MethodGet, c.tunnelUrl, nil)
	request.Host = host
	request.Header.Set(HEADER_REMOTE_ADDR, c.remoteAddr)
	request.Header.Set("HOST", request.Host)
	err = request.Write(tunnelConn)
	if err != nil {
		log.Error("send connect request err", err)
		return nil
	}
	// receive response
	response, err := http.ReadResponse(bufio.NewReader(tunnelConn), request)
	if err != nil {
		log.Error("receive connect response err", err)
		return nil
	}
	if response.StatusCode != http.StatusOK {
		log.Error("connect http tunnel err")
		return nil
	}
	return tunnelConn
}

type clientConn struct {
	net.Conn
}

func newClientConn(conn net.Conn) *clientConn {
	setTCPConnOptions(conn)
	cc := &clientConn{
		Conn: conn,
	}
	return cc
}
