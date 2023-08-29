package tunnel

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
	"net"
	"net/http"
)

type Client struct {
	config *ClientConfig
	// it is not nil if IsSmux is true
	smuxSession *smux.Session
}

func NewClient(cc *ClientConfig, options ...ClientOption) *Client {
	c := &Client{
		config: cc,
	}
	for _, option := range options {
		option(c)
	}
	log.Infof("NewClient config %s", c.config)
	return c
}

func (c *Client) ListenAndServe() error {
	l, err := net.Listen("tcp", c.config.LocalAddr)
	if err != nil {
		log.Errorf("listen LocalAddr %s, err: %v", c.config.LocalAddr, err)
		return err
	}
	defer l.Close()
	log.Infof("listen LocalAddr %s", c.config.LocalAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("accept, err: %v", err)
			continue
		}
		go c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn net.Conn) {
	log.Infof("handle conn %s", conn.RemoteAddr())
	defer conn.Close()
	setTCPConnOptions(conn)

	connectFn := func() net.Conn {
		var tunnelConn net.Conn
		switch c.config.Mode {
		case CONNECT_HTTP:
			tunnelConn = c.connectWithHTTP()
		case CONNECT_WEBSOCKET:
			tunnelConn = c.connectWithWebSocket()
		default:
			panic("mode is empty")
		}
		return tunnelConn
	}

	// support smux
	if c.config.IsSmux {
		if c.smuxSession == nil || c.smuxSession.IsClosed() {
			tunnelConn := connectFn()
			if tunnelConn == nil {
				return
			}
			defer tunnelConn.Close()
			session, err := smux.Client(tunnelConn, smux.DefaultConfig())
			if err != nil {
				log.Errorf("new IsSmux client, err: %v", err)
				return
			}
			defer session.Close()
			c.smuxSession = session
		}
		stream, err := c.smuxSession.OpenStream()
		if err != nil {
			log.Errorf("mux open stream, err: %v", err)
			return
		}
		if stream == nil {
			log.Errorf("mux open stream is null")
			return
		}
		copyDataOnConn(conn, stream)
		return
	}

	// per connection
	tunnelConn := connectFn()
	if tunnelConn == nil {
		return
	}
	defer tunnelConn.Close()
	copyDataOnConn(conn, tunnelConn)
}

func (c *Client) setHeader(header *http.Header) {
	header.Set(HEADER_MODE, string(c.config.Mode))
	header.Set(HEADER_REMOTE_ADDR, c.config.RemoteAddr)
	header.Set(HEADER_IS_SMUX, fmt.Sprint(c.config.IsSmux))
	header.Set(HEADER_TOKEN, c.config.Token)
}
