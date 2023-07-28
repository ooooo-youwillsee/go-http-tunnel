package tunnel

import (
	"errors"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
)

var (
	URL_CONNECT = "/__connect__"

	HEADER_REMOTE_ADDR = "REMOTE-ADDR"

	ErrAuthFail = errors.New("auth fail")
)

type Server struct {
	addr string
	url  string
	l    net.Listener
}

func NewServer(addr string, url string, options ...ServerOption) *Server {
	if url == "" {
		url = URL_CONNECT
	}
	s := &Server{
		addr: addr,
		url:  url,
	}
	for _, option := range options {
		option(s)
	}

	log.Infof("NewServer addr[%s], url[%s]", addr, url)
	return s
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Error("listen localAddr err ", err)
		return err
	}
	s.l = l
	log.Infof("listen localAddr %s", s.addr)

	mux := http.NewServeMux()
	mux.HandleFunc(s.url, s.Connect)
	err = http.Serve(s.l, mux)
	if err != nil {
		log.Error("serveHTTP err", err)
	}
	return err
}

func (s *Server) Connect(w http.ResponseWriter, r *http.Request) {
	upgrade := r.Header.Get("Upgrade")
	if upgrade == "websocket" {
		s.ConnectWithWebsocket(w, r)
		return
	}

	// auth client
	remoteAddr, err := s.auth(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

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
	c := s.newServerConn(conn)
	c.remoteAddr = remoteAddr
	go s.handleConn(c)
}

func (s *Server) ConnectWithWebsocket(w http.ResponseWriter, r *http.Request) {
	// upgrade http to websocket
	upgrade := websocket.Upgrader{}
	wsc, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer wsc.Close()
	// auth client
	remoteAddr, err := s.auth(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Error("dial remoteAddr err", err)
		return
	}
	defer remoteConn.Close()

	go func() {
		for {
			_, bytes, err := wsc.ReadMessage()
			if err != nil {
				log.Error(err)
				break
			}
			_, _ = remoteConn.Write(bytes)

			//_, r, err := wsc.NextReader()
			//if err != nil {
			//	log.Error(err)
			//	break
			//}
			//io.Copy(remoteConn, r)
		}
	}()

	for {
		buf := make([]byte, 1024)
		n, err := remoteConn.Read(buf)
		if err != nil {
			log.Error(err)
			break
		}
		_ = wsc.WriteMessage(websocket.BinaryMessage, buf[:n])

		//w, err := wsc.NextWriter(websocket.BinaryMessage)
		//if err != nil {
		//	log.Error(err)
		//	break
		//}
		//io.Copy(w, remoteConn)
		//w.Close()
	}
	//c := s.newServerConn(conn)
	//c.remoteAddr = remoteAddr
	//go s.handleConn(c)
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		log.Errorf("auth method '%s' is not supported", r.Method)
		return "", ErrAuthFail
	}

	// get remote localAddr
	remoteAddr := r.Header.Get(HEADER_REMOTE_ADDR)
	if remoteAddr == "" {
		log.Errorf("http header '%s' not found", HEADER_REMOTE_ADDR)
		return "", ErrAuthFail
	}
	return remoteAddr, nil
}

func (s *Server) handleConn(conn *ServerConn) {
	defer conn.Close()

	// copy data
	//errCh := Copy(conn, remoteConn)
	//<-errCh
}

type ServerConn struct {
	io.ReadWriteCloser
	remoteAddr string
}

func (s *Server) newServerConn(conn io.ReadWriteCloser) *ServerConn {
	c := &ServerConn{
		ReadWriteCloser: conn,
	}
	return c
}
