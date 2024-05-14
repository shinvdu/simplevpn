package client

import (
	"fmt"
	"log"
	"net"

	"github.com/shinvdu/simplevpn/common/cipher"
	"github.com/shinvdu/simplevpn/common/config"
	"github.com/shinvdu/simplevpn/vpn"

	// "github.com/songgao/water/waterutil"
	"github.com/songgao/packets/ethernet"
	"golang.org/x/net/ipv4"
)

func StartUDPClient(config config.Config) {
	iface := vpn.CreateVpn(config.CIDR)
	serverAddr, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	log.Printf("serverAddr: %s\n", serverAddr)

	if err != nil {
		log.Fatalln("failed to resolve server addr:", err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	log.Printf("localAddr: %s\n", localAddr)

	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)

	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}
	defer conn.Close()
	log.Printf("govpn udp client started on %v,CIDR is %v", config.LocalAddr, config.CIDR)

	go func() {
		buf := make([]byte, 1500)
		var frame ethernet.Frame
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil || n == 0 {
				continue
			}
			frame.Resize(n)
			b := cipher.XOR(buf[:n])
			copy(frame[:n], b[:n])
			// if !waterutil.IsIPv4(b) {
			// 	continue
			// }
			// log.Printf("Dst: %s\n", frame.Destination())
			// log.Printf("Src: %s\n", frame.Source())
			// log.Printf("Ethertype: % x\n", frame.Ethertype())
			ip_payload := frame.Payload()
			// log.Printf("Payload: % x\n", ip_payload)
			header, _ := ipv4.ParseHeader(ip_payload[:])
			// fmt.Printf("Received %d bytes from remote, header: %+v\n", n, header)
			log.Printf("read remote: %s => %s \n", header.Src, header.Dst)

			iface.Write(b)
		}
	}()

	// packet := make([]byte, 1500)
	var packet ethernet.Frame

	// conn.WriteToUDP(b, serverAddr)
	// create a byte slic with "hello"
	hello := []byte("hello")
	// get all data from slice hello

	// hello_enc := cipher.XOR(packet[:])
	conn.WriteToUDP(hello, serverAddr)

	for {
		packet.Resize(1500)
		// n, err := iface.Read(packet)
		n, err := iface.Read([]byte(packet))

		// fmt.Println("rechive data size:", n)
		// fmt.Println(packet)
		if err != nil || n == 0 {
			// fmt.Println("err iface data error:", err)
			continue
		}
		// if !waterutil.IsIPv4(packet) {
		// 	fmt.Println("not a IsIPv4 package, skip")
		// 	continue
		// }
		packet = packet[:n]

		b := cipher.XOR(packet[:n])

		// log.Printf("Dst: %s\n", packet.Destination())
		// log.Printf("Src: %s\n", packet.Source())
		// log.Printf("Ethertype: % x\n", packet.Ethertype())
		ip_payload := packet.Payload()
		// log.Printf("Payload: % x\n", ip_payload)
		header, _ := ipv4.ParseHeader(ip_payload[:])
		// get src from header

		// fmt.Printf("Received %d bytes from iface, header: %+v\n", n, header)
		log.Printf("write remote: %s => %s \n", header.Src, header.Dst)

		fmt.Printf("WriteToUDP: %v \n", serverAddr)
		conn.WriteToUDP(b, serverAddr)
	}
}
