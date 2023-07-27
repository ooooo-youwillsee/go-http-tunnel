package tunnel

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type Client struct {
	addr       string
	l          net.Listener
	remoteAddr string
	remoteConn net.Conn
}

func NewClient(addr, remoteAddr string) *Client {
	return &Client{
		addr:       addr,
		remoteAddr: remoteAddr,
	}
}

func (c *Client) ListenAndServe() error {
	l, err := net.Listen("tcp", c.addr)
	if err != nil {
		log.Fatalln("listen addr err", err)
		return err
	}
	c.l = l
	log.Printf("listen addr %s\n", c.addr)

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
	remoteConn, err := net.Dial("tcp", c.remoteAddr)
	if err != nil {
		log.Errorf("dial addr %s err %v\n", c.remoteAddr, err)
		return
	}
	c.remoteConn = remoteConn

	// send request
	request, _ := http.NewRequest(http.MethodConnect, URL_CONNECT, nil)
	err = request.Write(c.remoteConn)
	if err != nil {
		log.Error("send connect request err", err)
		return
	}
	// receive response
	response, err := http.ReadResponse(bufio.NewReader(c.remoteConn), request)
	if err != nil {
		log.Error("receive connect response err", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Error("connect http tunnel err")
		return
	}
	// copy data
	errCh := CopyConn(conn, remoteConn)
	<-errCh
}
