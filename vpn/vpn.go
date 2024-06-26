package vpn

import (
	"log"

	"github.com/shinvdu/simplevpn/common/osutil"
	"github.com/songgao/water"
)

func CreateVpn(cidr string) (iface *water.Interface) {
	c := water.Config{DeviceType: water.TAP}
	iface, err := water.New(c)
	if err != nil {
		log.Fatalln("failed to allocate vpn interface:", err)
	}
	log.Println("interface allocated:", iface.Name())
	osutil.ConfigVpn(cidr, iface)
	return iface
}

