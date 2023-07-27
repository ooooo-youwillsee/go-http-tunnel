package tunnel

type ServerOption func(server *Server)

type ClientOption func(client *Client)

func ServerWithTCPMapping(tcpProxies *TCPProxies) ServerOption {
	return func(server *Server) {
		server.tcpms = tcpProxies
	}
}

func ClientWithTCPProxies(tcpProxies *TCPProxies) ClientOption {
	return func(server *Client) {
		server.tcpProxies = tcpProxies
	}
}
