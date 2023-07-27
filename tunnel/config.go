package tunnel

import (
	"github.com/spf13/viper"
	"strings"
)

type TCPProxies []*TCPProxy

type TCPProxy struct {
	LocalAddr  string
	RemoteAddr string
	TunnelAddr string
}

func NewTcpProxiesFromFile(configPath string) *TCPProxies {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("read config err for path " + configPath)
	}

	tcpProxies := &TCPProxies{}
	var localAddr string
	var remoteAddr string
	var tunnelAddr string

	for _, m := range viper.AllSettings() {
		for k, v := range m.(map[string]interface{}) {
			switch k {
			case "local_addr":
				localAddr = v.(string)
			case "remote_addr":
				remoteAddr = v.(string)
			case "tunnel_addr":
				tunnelAddr = v.(string)
			}
		}
		tcpProxies.AddTCPProxy(localAddr, remoteAddr, tunnelAddr)
	}
	return tcpProxies
}

func (ms *TCPProxies) AddTCPProxy(localAddr, remoteAddr, tunnelAddr string) {
	m := &TCPProxy{
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
		TunnelAddr: tunnelAddr,
	}
	*ms = append(*ms, m)
}

func (ms *TCPProxies) GetTCPProxy(localAddr string) *TCPProxy {
	for _, m := range *ms {
		_, port1 := splitAddr(m.LocalAddr)
		_, port2 := splitAddr(localAddr)
		if port1 == port2 {
			return m
		}
	}
	return nil
}

func (ms *TCPProxies) GetTunnelAddrs() []string {
	tunnelAddrs := make([]string, 0)
	m := make(map[string]struct{})
	for _, proxy := range *ms {
		tunnelAddr := proxy.TunnelAddr
		if _, ok := m[tunnelAddr]; !ok {
			m[tunnelAddr] = struct{}{}
			tunnelAddrs = append(tunnelAddrs, tunnelAddr)
		}
	}
	return tunnelAddrs
}

func splitAddr(addr string) (string, string) {
	split := strings.Split(addr, ":")
	return split[0], split[1]
}
