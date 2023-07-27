package tunnel

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
)

var (
	URL_CONNECT = "/__connect__"
)

type Server struct {
	addr string
	l    net.Listener
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Error("listen addr err ", err)
		return err
	}
	s.l = l
	log.Printf("listen addr %s\n", s.addr)

	mux := http.NewServeMux()
	mux.HandleFunc(URL_CONNECT, s.Connect)
	err = http.Serve(s.l, mux)
	if err != nil {
		log.Error("serveHTTP err", err)
	}
	return err
}

func (s *Server) Connect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "connected success")

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "not support hijacker")
		return
	}

	conn, _, _ := hijacker.Hijack()
	c := newServerConn(conn)
	go s.handleConn(c)
}

func (s *Server) handleConn(conn *ServerConn) {
	defer conn.Close()
	remoteConn, err := net.Dial("tcp", conn.remoteAddr)
	if err != nil {
		log.Error("dial addr err", err)
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

func newServerConn(conn net.Conn) *ServerConn {
	return &ServerConn{
		Conn:       conn,
		remoteAddr: ":30001",
	}
}
