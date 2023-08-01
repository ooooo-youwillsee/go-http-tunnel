package tunnel

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

var (
	URL_CONNECT = "/"

	HEADER_REMOTE_ADDR = "REMOTE-ADDR"
	HEADER_TOKEN       = "TOKEN"

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
	remoteAddr, err := s.auth(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	upgrade := r.Header.Get("Upgrade")
	if upgrade == "websocket" {
		s.connectWithWebsocket(w, r, remoteAddr)
		return
	}
	s.connectWithHTTP(w, r, remoteAddr)
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
	token := r.Header.Get(HEADER_TOKEN)
	if token != s.token {
		log.Errorf("http header token '%s' is err", token)
		return "", ErrAuthFail
	}
	return remoteAddr, nil
}
