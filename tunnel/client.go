package tunnel

import (
	log "github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
	"net"
)

var (
	CONNECT_HTTP      ConnectMode = "http"
	CONNECT_WEBSOCKET ConnectMode = "websocket"
)

type ConnectMode string

type Client struct {
	localAddr  string
	remoteAddr string
	tunnelAddr string
	tunnelUrl  string
	token      string
	mode       ConnectMode
	// true or false,  default is true
	isSmux      string
	smuxSession *smux.Session
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

	connectFn := func() net.Conn {
		var tunnelConn net.Conn
		switch c.mode {
		case CONNECT_HTTP:
			tunnelConn = c.connectWithHTTP()
		case CONNECT_WEBSOCKET:
			tunnelConn = c.connectWithWebSocket()
		}
		return tunnelConn
	}

	// support isSmux
	if c.isSmux == "true" {
		if c.smuxSession == nil || c.smuxSession.IsClosed() {
			tunnelConn := connectFn()
			defer tunnelConn.Close()
			session, err := smux.Client(tunnelConn, smux.DefaultConfig())
			if err != nil {
				log.Error("new isSmux client ", err)
				return
			}
			defer session.Close()
			c.smuxSession = session
		}
		stream, err := c.smuxSession.OpenStream()
		if err != nil {
			log.Error("mux open stream ", err)
			return
		}
		copyDataOnConn(conn, stream)
		return
	}

	// per connection
	tunnelConn := connectFn()
	defer tunnelConn.Close()
	copyDataOnConn(conn, tunnelConn)
}
