package main

import (
	"fmt"
	"github.com/liuxinglanyue/mtr"
	"net"
	"os"
	"time"
)

const default_local_addr = "0.0.0.0"

func RealTimePing(dest string) (hop mtr.TracerouteReturn, err error) {
	ipAddr := net.IPAddr{IP: net.ParseIP(dest)}
	pid := os.Getpid() & 0xffff
	ttl := 64
	timeout := time.Duration(2000) * time.Millisecond
	hop, err = mtr.Icmp(default_local_addr, &ipAddr, ttl, pid, timeout)
	return hop, err
}

func main() {
	hop, err := RealTimePing("www.baidu.com")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(hop.Success, hop.Addr, hop.Elapsed)
}
