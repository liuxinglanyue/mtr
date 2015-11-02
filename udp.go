package mtr

import (
	"syscall"
	"time"
)

func Udp(socketAddr, destAddr [4]byte, ttl, port int, tv syscall.Timeval, p []byte) (hop TracerouteReturn, err error) {
	start := time.Now()

	recvSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return hop, err
	}

	sendSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return hop, err
	}
	syscall.SetsockoptInt(sendSocket, 0x0, syscall.IP_TTL, ttl)
	syscall.SetsockoptTimeval(recvSocket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
	defer syscall.Close(recvSocket)
	defer syscall.Close(sendSocket)

	syscall.Bind(recvSocket, &syscall.SockaddrInet4{Port: port, Addr: socketAddr})
	syscall.Sendto(sendSocket, []byte{0x0}, 0, &syscall.SockaddrInet4{Port: port, Addr: destAddr})

	_, from, err := syscall.Recvfrom(recvSocket, p, 0)
	elapsed := time.Since(start)
	if err == nil {
		currAddr := from.(*syscall.SockaddrInet4).Addr
		hop.Addr = AddressString(currAddr)
		hop.Success = true
		hop.Elapsed = elapsed
	}

	return hop, err
}
