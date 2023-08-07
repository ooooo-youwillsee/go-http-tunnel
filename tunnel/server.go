package tunnel

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
	"net"
	"net/http"
)

var (
	URL_CONNECT = "/"

	HEADER_REMOTE_ADDR = "REMOTE-ADDR"
	HEADER_TOKEN       = "TOKEN"
	HEADER_IS_SMUX     = "IS_SMUX"

	ErrAuthFail = errors.New("auth fail")
)

type Server struct {
	addr  string
	url   string
	token string
	l     net.Listener
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
	// auth client
	remoteConn, err := s.auth(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer remoteConn.Close()

	var conn net.Conn
	upgrade := r.Header.Get("Upgrade")
	if upgrade == "websocket" {
		conn = s.connectWithWebSocket(w, r)
	} else {
		conn = s.connectWithHTTP(w, r)
	}
	defer conn.Close()

	// support isSmux
	isSmux := r.Header.Get(HEADER_IS_SMUX)
	if isSmux == "true" {
		session, err := smux.Server(conn, smux.DefaultConfig())
		if err != nil {
			log.Error("new isSmux client ", err)
			return
		}
		defer session.Close()

		for {
			stream, err := session.AcceptStream()
			if err != nil {
				log.Error("isSmux open stream ", err)
				return
			}
			go copyDataOnConn(stream, remoteConn)
		}
	}

	// per connection
	copyDataOnConn(conn, remoteConn)
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	// get remote localAddr
	remoteAddr := r.Header.Get(HEADER_REMOTE_ADDR)
	if remoteAddr == "" {
		log.Errorf("http header '%s' not found", HEADER_REMOTE_ADDR)
		return nil, ErrAuthFail
	}
	// verify token
	token := r.Header.Get(HEADER_TOKEN)
	if token != s.token {
		log.Errorf("http header token '%s' is err", token)
		return nil, ErrAuthFail
	}
	// dial remote addr
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Error("dial remote addr ", err)
		return nil, err
	}
	return remoteConn, nil
}
