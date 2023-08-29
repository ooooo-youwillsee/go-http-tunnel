package tunnel

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
)

const (
	DEFAULT_TUNNEL_URL = "/"

	CONFIG_COMMON      = "common"
	CONFIG_LOCAL_ADDR  = "local_addr"
	CONFIG_REMOTE_ADDR = "remote_addr"
	CONFIG_TUNNEL_ADDR = "tunnel_addr"
	CONFIG_TUNNEL_URL  = "tunnel_url"
	CONFIG_TOKEN       = "Token"
	CONFIG_SMUX        = "IsSmux"
	CONFIG_MODE        = "Mode"
)

const (
	CONNECT_HTTP      ConnectMode = "http"
	CONNECT_WEBSOCKET ConnectMode = "websocket"
)

type ConnectMode string

type ClientConfigs []*ClientConfig

type ClientConfig struct {
	LocalAddr  string
	RemoteAddr string
	TunnelAddr string
	TunnelUrl  string
	Token      string
	IsSmux     bool
	Mode       ConnectMode
}

func NewClientConfigsFromFile(configFile string) *ClientConfigs {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("read config err for path " + configFile)
	}

	ccs := &ClientConfigs{}
	// parse common
	tunnelAddr := viper.GetString(CONFIG_COMMON + "." + CONFIG_TUNNEL_ADDR)
	tunnelUrl := viper.GetString(CONFIG_COMMON + "." + CONFIG_TUNNEL_URL)
	isSmux := viper.GetBool(CONFIG_COMMON + "." + CONFIG_SMUX)
	mode := viper.GetString(CONFIG_COMMON + "." + CONFIG_MODE)
	if tunnelUrl == "" {
		tunnelUrl = DEFAULT_TUNNEL_URL
	}
	if mode == "" {
		mode = string(CONNECT_WEBSOCKET)
	}

	// parse special
	for g, m := range viper.AllSettings() {
		if g == CONFIG_COMMON {
			continue
		}
		cc := &ClientConfig{}
		cc.LocalAddr = getValue(m, CONFIG_LOCAL_ADDR, "")
		cc.RemoteAddr = getValue(m, CONFIG_REMOTE_ADDR, "")
		cc.TunnelAddr = getValue(m, CONFIG_TUNNEL_ADDR, tunnelAddr)
		cc.TunnelUrl = getValue(m, CONFIG_TUNNEL_URL, tunnelUrl)
		cc.Token = getValue(m, CONFIG_TOKEN, "")
		cc.IsSmux = getValue(m, CONFIG_SMUX, isSmux)
		cc.Mode = ConnectMode(getValue(m, CONFIG_MODE, mode))
		if cc.LocalAddr == "" || cc.RemoteAddr == "" || cc.TunnelAddr == "" {
			panic(fmt.Sprintf("group %s LocalAddr or RemoteAddr or TunnelAddr is empty", g))
		}
		*ccs = append(*ccs, cc)
	}
	return ccs
}

func (c *ClientConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

type ServerConfigs []*ServerConfig

type ServerConfig struct {
	TunnelAddr string
	TunnelUrl  string
	Token      string
}

func NewServerConfigsFrom(configFile string) *ServerConfigs {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("read config err for path " + configFile)
	}

	scs := &ServerConfigs{}
	for g, m := range viper.AllSettings() {
		sc := &ServerConfig{}
		sc.TunnelAddr = getValue(m, CONFIG_TUNNEL_ADDR, "")
		sc.TunnelUrl = getValue(m, CONFIG_TUNNEL_URL, DEFAULT_TUNNEL_URL)
		sc.Token = getValue(m, CONFIG_TOKEN, "")
		if sc.TunnelAddr == "" {
			panic(fmt.Sprintf("group %s TunnelAddr is empty", g))
		}
		*scs = append(*scs, sc)
	}
	return scs
}

func (c *ServerConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func getValue[T any](m interface{}, key string, defaultValue T) T {
	mm := m.(map[string]interface{})
	if v, ok := mm[key]; ok {
		return v.(T)
	}
	return defaultValue
}
