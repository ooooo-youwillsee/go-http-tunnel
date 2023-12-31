package tunnel

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) connectWithWebSocket() net.Conn {
	wsurl := url.URL{
		Scheme: "ws",
		Host:   c.config.TunnelAddr,
		Path:   c.config.TunnelUrl,
	}
	header := http.Header{}
	c.setHeader(&header)
	wsc, _, err := websocket.DefaultDialer.Dial(wsurl.String(), header)
	if err != nil {
		log.Errorf("dial websocket addr %s, err: %v", wsurl.String(), err)
		return nil
	}
	conn := newWebSocketConn(wsc)
	return conn
}

func (s *Server) connectWithWebSocket(w http.ResponseWriter, r *http.Request) net.Conn {
	// upgrade http to websocket
	upgrader := websocket.Upgrader{}
	wsc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade http to websocket, err: %v", err)
		return nil
	}
	conn := newWebSocketConn(wsc)
	return conn
}

type wsConn struct {
	wsc *websocket.Conn
	buf []byte
}

func newWebSocketConn(wsc *websocket.Conn) net.Conn {
	return &wsConn{
		wsc: wsc,
	}
}

func (w *wsConn) Read(b []byte) (n int, err error) {
	if len(w.buf) > 0 {
		n = copy(b, w.buf)
		w.buf = w.buf[n:]
		return
	}
	_, buf, err := w.wsc.ReadMessage()
	if err != nil {
		return 0, err
	}
	n = copy(b, buf)
	w.buf = buf[n:]
	return
}

func (w *wsConn) Write(b []byte) (n int, err error) {
	err = w.wsc.WriteMessage(websocket.BinaryMessage, b)
	n = len(b)
	return
}

func (w *wsConn) Close() error {
	return w.wsc.Close()
}

func (w *wsConn) LocalAddr() net.Addr {
	return w.wsc.LocalAddr()
}

func (w *wsConn) RemoteAddr() net.Addr {
	return w.wsc.RemoteAddr()
}

func (w *wsConn) SetDeadline(t time.Time) error {
	return w.wsc.UnderlyingConn().SetDeadline(t)
}

func (w *wsConn) SetReadDeadline(t time.Time) error {
	return w.wsc.UnderlyingConn().SetReadDeadline(t)
}

func (w *wsConn) SetWriteDeadline(t time.Time) error {
	return w.wsc.SetWriteDeadline(t)
}
