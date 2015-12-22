package mtr

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"syscall"
	"time"
)

func Icmp(localAddr string, dst net.Addr, ttl, pid int, timeout time.Duration) (hop TracerouteReturn, err error) {
	hop.Success = false
	start := time.Now()
	c, err := icmp.ListenPacket("ip4:icmp", localAddr)
	if err != nil {
		return hop, err
	}
	defer c.Close()
	c.IPv4PacketConn().SetTTL(ttl)
	c.SetDeadline(time.Now().Add(timeout))
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: pid, Seq: 1,
			Data: []byte(""),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return hop, err
	}

	if _, err := c.WriteTo(wb, dst); err != nil {
		return hop, err
	}

	rb := make([]byte, 1500)
	_, peer, err := c.ReadFrom(rb)
	if err != nil {
		return hop, err
	}
	elapsed := time.Since(start)
	hop.Elapsed = elapsed
	hop.Addr = peer.String()
	hop.Success = true
	return hop, err
}

func IcmpRpc(localAddr string, dstAddr string, ttl int, timeout int) (hop TracerouteReturn, err error) {
	ipAddr := net.IPAddr{IP: net.ParseIP(dstAddr)}
	pid := os.Getpid() & 0xffff
	to := time.Duration(timeout) * time.Millisecond
	hop, err = Icmp(localAddr, &ipAddr, ttl, pid, to)
	return hop, err
}

func IcmpWrapper(localAddr, destAddr [4]byte, ttl, port int, tv syscall.Timeval, p []byte) (hop TracerouteReturn, err error) {
	ipAddr := net.IPAddr{IP: net.ParseIP(AddressString(destAddr))}
	timeout := time.Duration(tv.Nano()/(1000*1000)) * time.Millisecond
	pid := os.Getpid() & 0xffff
	return Icmp(AddressString(localAddr), &ipAddr, ttl, pid, timeout)
}
