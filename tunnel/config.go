package tunnel

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type ClientConfigs []*ClientConfig

type ClientConfig struct {
	LocalAddr  string
	RemoteAddr string
	TunnelAddr string
	TunnelUrl  string
	Token      string
	Smux       string
}

func NewClientConfigsFromFile(configFile string) *ClientConfigs {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("read config err for path " + configFile)
	}

	ccs := &ClientConfigs{}
	// parse common
	tunnelAddr := viper.GetString("common.tunnel_addr")
	tunnelUrl := viper.GetString("common.tunnel_url")

	// parse special
	for g, m := range viper.AllSettings() {
		if g == "common" {
			continue
		}
		cc := &ClientConfig{}
		cc.LocalAddr = getString(m, "local_addr", "")
		cc.RemoteAddr = getString(m, "remote_addr", "")
		cc.TunnelAddr = getString(m, "tunnel_addr", tunnelAddr)
		cc.TunnelUrl = getString(m, "tunnel_url", tunnelUrl)
		cc.Token = getString(m, "token", "")
		cc.Smux = getString(m, "isSmux", "true")
		*ccs = append(*ccs, cc)
	}
	return ccs
}

type ServerConfigs []*ServerConfig

type ServerConfig struct {
	Addr  string
	Url   string
	Token string
}

func NewServerConfigsFrom(configFile string) *ServerConfigs {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("read config err for path " + configFile)
	}

	scs := &ServerConfigs{}
	for g, m := range viper.AllSettings() {
		sc := &ServerConfig{}
		sc.Addr = getString(m, "tunnel_addr", "")
		sc.Url = getString(m, "tunnel_url", "")
		sc.Token = getString(m, "token", "")
		if sc.Addr == "" {
			panic(fmt.Sprintf("group %s config 'Addr' is empty", g))
		}
		*scs = append(*scs, sc)
	}
	return scs
}

func splitAddr(addr string) (string, string) {
	split := strings.Split(addr, ":")
	return split[0], split[1]
}

func getString(m interface{}, key string, defaultValue string) string {
	mm := m.(map[string]interface{})
	if v, ok := mm[key]; ok {
		return v.(string)
	}
	return defaultValue
}
