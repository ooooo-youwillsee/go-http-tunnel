package tunnel

import (
	"errors"
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
	addr  string
	l     net.Listener
	tcpms *TCPProxies
}

func NewServer(addr string, options ...ServerOption) *Server {
	s := &Server{
		addr: addr,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Error("listen localAddr err ", err)
		return err
	}
	s.l = l
	log.Printf("listen localAddr %s\n", s.addr)

	mux := http.NewServeMux()
	mux.HandleFunc(URL_CONNECT, s.Connect)
	err = http.Serve(s.l, mux)
	if err != nil {
		log.Error("serveHTTP err", err)
	}
	return err
}

func (s *Server) Connect(w http.ResponseWriter, r *http.Request) {
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

	conn, _, _ := hijacker.Hijack()
	c := s.newServerConn(conn)
	c.remoteAddr = remoteAddr
	go s.handleConn(c)
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodConnect {
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
	remoteConn, err := net.Dial("tcp", conn.remoteAddr)
	if err != nil {
		log.Error("dial remoteAddr err", err)
		return
	}
	defer remoteConn.Close()
	// copy data
	errCh := CopyConn(conn, remoteConn)
	<-errCh
}

type ServerConn struct {
	net.Conn
	remoteAddr string
}

func (s *Server) newServerConn(conn net.Conn) *ServerConn {
	c := &ServerConn{
		Conn: conn,
	}
	return c
}
