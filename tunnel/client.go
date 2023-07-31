package tunnel

import (
	log "github.com/sirupsen/logrus"
	"net"
)

var (
	CONNECT_HTTP      connectMode = "http"
	CONNECT_WEBSOCKET connectMode = "websocket"
)

type connectMode string

type Client struct {
	localAddr  string
	remoteAddr string
	tunnelAddr string
	tunnelUrl  string
	mode       connectMode
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
		mode:       CONNECT_WEBSOCKET,
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
		go c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn net.Conn) {
	defer conn.Close()
	setTCPConnOptions(conn)
	// connect
	switch c.mode {
	case CONNECT_HTTP:
		c.connectWithHTTP(conn)
	case CONNECT_WEBSOCKET:
		c.connectWithWebsocket(conn)
	}
}
