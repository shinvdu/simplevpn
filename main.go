package main

import (
	"flag"
	"fmt"

	"github.com/shinvdu/simplevpn/client"
	"github.com/shinvdu/simplevpn/common/config"
	"github.com/shinvdu/simplevpn/server"
)

func main() {
	config := config.Config{}
	flag.StringVar(&config.CIDR, "c", "172.16.0.1/24", "vpn interface CIDR")
	flag.StringVar(&config.LocalAddr, "l", "0.0.0.0:3000", "local address")
	flag.StringVar(&config.ServerAddr, "s", "0.0.0.0:3001", "server address")
	flag.StringVar(&config.Key, "k", "6w9z$C&F)J@NcRfWjXn3r4u7x!A%D*G-", "encryption key")
	flag.StringVar(&config.Protocol, "p", "udp", "protocol udp")
	flag.BoolVar(&config.ServerMode, "S", false, "server mode")
	flag.Parse()
	config.Init()
	fmt.Println(config)
	switch config.Protocol {
	case "udp":
		if config.ServerMode {
			server.StartUDPServer(config)
		} else {
			client.StartUDPClient(config)
		}
	default:
	}
}
