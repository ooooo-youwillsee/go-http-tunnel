package tunnel

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/url"
	"sync"
)

func (c *Client) connectWithWebsocket(conn net.Conn) {
	wsurl := url.URL{
		Scheme: "ws",
		Host:   c.tunnelAddr,
		Path:   c.tunnelUrl,
	}
	header := http.Header{}
	header.Set(HEADER_REMOTE_ADDR, c.remoteAddr)
	header.Set(HEADER_TOKEN, c.token)
	wsc, _, err := websocket.DefaultDialer.Dial(wsurl.String(), header)
	if err != nil {
		log.Error("dial websocket addr ", err)
		return
	}
	defer wsc.Close()

	var wg sync.WaitGroup
	// read data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConnToWebsocket(conn, wsc)
	}()

	// write data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyWebsocketToConn(wsc, conn)
	}()
	wg.Wait()
}

func (s *Server) connectWithWebsocket(w http.ResponseWriter, r *http.Request, remoteAddr string) {
	// upgrade http to websocket
	upgrader := websocket.Upgrader{}
	wsc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade http to websocket ", err)
		return
	}
	defer wsc.Close()

	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Error("dial remoteAddr ", err)
		return
	}
	defer remoteConn.Close()

	var wg sync.WaitGroup
	// read data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyWebsocketToConn(wsc, remoteConn)
	}()

	// write data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConnToWebsocket(remoteConn, wsc)
	}()
	wg.Wait()
}

func copyConnToWebsocket(conn net.Conn, wsc *websocket.Conn) {
	buf := make([]byte, 1024*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Error("conn read ", err)
			return
		}
		if n == 0 {
			continue
		}
		err = wsc.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			log.Error("websocket write ", err)
			return
		}
	}
}

func copyWebsocketToConn(wsc *websocket.Conn, conn net.Conn) {
	for {
		messageType, buf, err := wsc.ReadMessage()
		if err != nil {
			log.Error("websocket read ", err)
			return
		}
		if messageType != websocket.BinaryMessage {
			continue
		}
		_, err = conn.Write(buf)
		if err != nil {
			log.Error("conn write ", err)
			return
		}
	}
}
