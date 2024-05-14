package osutil

// imports
import (
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"fmt"

	"github.com/songgao/water"
)

func ConfigVpn(cidr string, iface *water.Interface) {
	os := runtime.GOOS
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}
	fmt.Println(os)
	if os == "linux" {
		// execCmd("/sbin/ip", "link", "del", "dev", iface.Name())

    // execCmd("/sbin/ip", "tuntap", "add", "dev", iface.Name(), "mode", "tap")
    // log.Printf("done for tuntap")

		execCmd("/sbin/ip", "addr", "add", cidr, "dev", iface.Name())
    log.Printf("done for addr")

		execCmd("/sbin/ip", "link", "set", "up", "dev", iface.Name())
    log.Printf("done for link")

	} else if os == "darwin" {
		minIp := ipNet.IP.To4()
		minIp[3]++
		execCmd("ifconfig", iface.Name(), "inet", ip.String(), minIp.String(), "up")
	} else if os == "windows" {
		log.Printf("please install openvpn client,see this link:%v", "https://github.com/OpenVPN/openvpn")
		log.Printf("open new cmd and enter:netsh interface ip set address name=\"%v\" source=static addr=%v mask=%v gateway=none", iface.Name(), ip.String(), ipNet.Mask.String())
	} else {
		log.Printf("not support os:%v", os)
	}
}

func execCmd(c string, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatalln("failed to exec /sbin/ip error:", err)
	}
}
