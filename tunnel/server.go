package tunnel

import (
	log "github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
	"net"
	"net/http"
	"strconv"
)

var (
	HEADER_REMOTE_ADDR = "REMOTE-ADDR"
	HEADER_TOKEN       = "TOKEN"
	HEADER_IS_SMUX     = "IS_SMUX"
	HEADER_MODE        = "mode"
)

type Server struct {
	config *ServerConfig
	l      net.Listener
}

func NewServer(sc *ServerConfig, options ...ServerOption) *Server {
	s := &Server{
		config: sc,
	}
	for _, option := range options {
		option(s)
	}
	log.Infof("NewServer config %s", s.config)
	return s
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.config.TunnelAddr)
	if err != nil {
		log.Error("listen LocalAddr %s, err: %v ", s.config.TunnelAddr, err)
		return err
	}
	s.l = l
	log.Infof("listen LocalAddr %s", s.config.TunnelAddr)

	mux := http.NewServeMux()
	mux.HandleFunc(s.config.TunnelUrl, s.Connect)
	err = http.Serve(s.l, mux)
	if err != nil {
		log.Errorf("serveHTTP err: %v", err)
	}
	return err
}

func (s *Server) Connect(w http.ResponseWriter, r *http.Request) {
	// auth client
	remoteAddr := s.auth(w, r)
	if remoteAddr == "" {
		http.NotFound(w, r)
		return
	}

	connectFn := func() net.Conn {
		var conn net.Conn
		mode := ConnectMode(r.Header.Get(HEADER_MODE))
		switch mode {
		case CONNECT_HTTP:
			conn = s.connectWithHTTP(w, r)
		case CONNECT_WEBSOCKET:
			conn = s.connectWithWebSocket(w, r)
		default:
			panic("header mode is empty")
		}
		return conn
	}

	remoteConnectFn := func() net.Conn {
		remoteConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			log.Errorf("dial remote addr %s, err: %v", remoteAddr, err)
			return nil
		}
		return remoteConn
	}

	// connect
	conn := connectFn()
	if conn == nil {
		return
	}
	defer conn.Close()

	// support smux
	isSmux, _ := strconv.ParseBool(r.Header.Get(HEADER_IS_SMUX))
	if isSmux {
		session, err := smux.Server(conn, smux.DefaultConfig())
		if err != nil {
			log.Errorf("new smux server, err: %v", err)
			return
		}
		defer session.Close()

		for {
			stream, err := session.AcceptStream()
			if err != nil {
				log.Errorf("smux accept stream, err: %v", err)
				return
			}
			if stream == nil {
				log.Errorf("smux accept stream is null")
				return
			}
			remoteConn := remoteConnectFn()
			if remoteConn == nil {
				continue
			}
			go copyDataOnConn(stream, remoteConn)
		}
	}

	// per connection
	remoteConn := remoteConnectFn()
	if remoteConn == nil {
		return
	}
	copyDataOnConn(conn, remoteConn)
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) string {
	// verify Token
	token := r.Header.Get(HEADER_TOKEN)
	if token != s.config.Token {
		log.Errorf("http header Token '%s' is err", token)
		return ""
	}
	// get remote LocalAddr
	remoteAddr := r.Header.Get(HEADER_REMOTE_ADDR)
	if remoteAddr == "" {
		log.Errorf("http header '%s' not found", HEADER_REMOTE_ADDR)
	}
	return remoteAddr
}
