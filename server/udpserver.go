package server

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/krishpranav/govpn/common/cipher"
	"github.com/krishpranav/govpn/common/config"
	"github.com/krishpranav/govpn/vpn"
	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
)

func StartUDPServer(config config.Config) {
	iface := vpn.CreateVpn(config.CIDR)
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}
	defer conn.Close()
	log.Printf("govpn udp server started on %v,CIDR is %v", config.LocalAddr, config.CIDR)

	forwarder := &Forwarder{localConn: conn, cliAddr: nil}
	go forwarder.forward(iface, conn)

	buf := make([]byte, 1500)
	hello := []byte("hello")
  var frame ethernet.Frame

	for {
		frame.Resize(1500)
		n, cliAddr, err := conn.ReadFromUDP(buf)
		// n, cliAddr, err := conn.ReadFromUDP([]byte(buf))
		fmt.Printf("rechive data size: %d, from : %v \n", n, cliAddr)
		// fmt.Printf("cliAddr: %v", cliAddr)

		if err != nil || n == 0 {
			continue
		}
		// 判断 hello 是否等于 buf[:5]

		// str := fmt.Sprint(slice)
		// fmt.Printf("first 5 chars: %s", buf[:5])
		// fmt.Printf("first 5 chars: %s", hello)

		if n > 4 && bytes.Equal(hello, buf[:5]) {
			// convert slice to string

			forwarder.cliAddr = cliAddr
			fmt.Printf("---------------- \n")
			fmt.Printf("---------------- \n")
			fmt.Printf("---------------- \n")
			fmt.Printf("get hello, seting peer: %v \n", cliAddr)
			fmt.Printf("---------------- \n")
			fmt.Printf("---------------- \n")
			fmt.Printf("---------------- \n")
			continue
		}

    b := cipher.XOR(buf[:n])
    frame = b[:n]

    ip_payload := frame.Payload()
    // log.Printf("Payload: % x\n", ip_payload)
    header, _ := ipv4.ParseHeader(ip_payload[:])
    // fmt.Printf("Received %d bytes from remtve, header: %+v\n", n, header)
    log.Printf("read remote: %s => %s \n", header.Src, header.Dst)
    
 
		// if !waterutil.IsIPv4(b) {
		// 	continue
		// }

		// log.Printf("read remote Dst: %s\n", buf.Destination())
		// log.Printf("read remote Src: %s\n", buf.Source())
		// log.Printf("read remote: %s => %s \n", buf.Source(), buf.Destination())
		// log.Printf("read remote Ethertype: % x\n", buf.Ethertype())

		iface.Write(b)
		// srcAddr, dstAddr := netutil.GetAddr(b)
		// if srcAddr == "" || dstAddr == "" {
		// 	continue
		// }
		// key := fmt.Sprintf("%v->%v", srcAddr, dstAddr)
		// fmt.Println(key)
		// if forwarder.cliAddr == nil {
		// 	forwarder.cliAddr = cliAddr
		// }
		// forwarder.connCache.Set(key, cliAddr, cache.DefaultExpiration)
	}
}

type Forwarder struct {
	localConn *net.UDPConn
	// cliAddr   net.Addr
  cliAddr *net.UDPAddr

	// connCache *cache.Cache
}

func (f *Forwarder) forward(iface *water.Interface, conn *net.UDPConn) {
	// packet := make([]byte, 1500)

	var frame ethernet.Frame
	for {
		frame.Resize(1500)
		// n, err := iface.Read(packet)
		n, err := iface.Read([]byte(frame))

		frame = frame[:n]

		// fmt.Println(packet)
		// fmt.Println(frame)

		if err != nil || n == 0 {
			continue
		}
		// log.Printf("read iface: %s => %s \n", buf.Source(), buf.Destination())

		// log.Printf("read iface Dst: %s\n", frame.Destination())
		// log.Printf("read iface Src: %s\n", frame.Source())
		// log.Printf("read iface Ethertype: % x\n", frame.Ethertype())
		ip_payload := frame.Payload()
		// log.Printf("Payload: % x\n", ip_payload)
		header, _ := ipv4.ParseHeader(ip_payload[:])
		// fmt.Printf("Received %d bytes from iface, header: %+v\n", n, header)
		log.Printf("read iface: %s => %s \n", header.Src, header.Dst)

		b := frame[:n]
		// if !waterutil.IsIPv4(b) {
		// 	continue
		// }
		// srcAddr, dstAddr := netutil.GetAddr(b)

		// if srcAddr == "" || dstAddr == "" {
		// 	log.Printf("skip, due to empty srcAddr and dstAddr: %v->%v", dstAddr, srcAddr)
		// 	continue
		// }
		// key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
		// log.Printf("srcAddr -> dstAddr: %v->%v", dstAddr, srcAddr)
		// fmt.Println(key)

		// v, ok := f.connCache.Get(key)
		// if ok {
		fmt.Printf("cliAddr: %v \n", f.cliAddr)

		b = cipher.XOR(b)

		if f.cliAddr == nil {
			fmt.Println("f.cliAdd empty, skip")
			continue
		}
		// cliUDPAddr := f.cliAddr.(*net.UDPAddr)
		// convert f.cliAddr into string
		// cliUDPAddrStr := &f.cliAddr

		// fmt.Printf("cliUDPAddrStr: %v \n", cliUDPAddrStr)
		// if cliUDPAddr != nil {
		// 	fmt.Println("cliUDPAddr, skip")
		// 	continue
		// }
		fmt.Printf("send back client cliAddr: %v \n", f.cliAddr)
		f.localConn.WriteToUDP(b, f.cliAddr)
	}
}
