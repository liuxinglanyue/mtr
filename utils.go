package mtr

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func LocalAddr() (addr [4]byte, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if len(ipnet.IP.To4()) == net.IPv4len {
				copy(addr[:], ipnet.IP.To4())
				return
			}
		}
	}
	err = errors.New("You do not appear to be connected to the Internet")
	return
}

func AddressString(addr [4]byte) string {
	return fmt.Sprintf("%v.%v.%v.%v", addr[0], addr[1], addr[2], addr[3])
}

func DestAddr(dest string) (destAddr [4]byte, err error) {
	addrs, err := net.LookupHost(dest)
	if err != nil {
		return
	}

	addr := addrs[0]

	ipAddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return
	}
	copy(destAddr[:], ipAddr.IP.To4())
	return
}

func DestAddrs(dest string) (addrs [][4]byte, err error) {
	addrStrs, err := net.LookupHost(dest)
	if err != nil {
		return
	}
	addrs = make([][4]byte, len(addrStrs))
	for i, addr := range addrStrs {
		ipAddr, err := net.ResolveIPAddr("ip", addr)
		if err != nil {
			continue
		}
		slice := [4]byte{}
		copy(slice[:], ipAddr.IP.To4())
		addrs[i] = slice
	}
	return
}

func Time2Float(t time.Duration) float32 {
	return (float32)(t/time.Microsecond) / float32(1000)
}

/*
func DestAddrOne(dest string) (destAddr [4]byte, err error) {
	destAddrs, err := DestAddr(dest)
	return destAddrs[0], err
}*/
