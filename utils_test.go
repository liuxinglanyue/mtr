package mtr

import (
	"testing"
	"time"
)

func TestLocalAddr(t *testing.T) {
	local, err := LocalAddr()
	if err != nil {
		t.Error(err)
	}
	t.Log(local[0], local[1], local[2], local[3])
}

func TestAddressString(t *testing.T) {
	ipByte := [4]byte{192, 168, 1, 23}
	ip := AddressString(ipByte)
	if !(ip == "192.168.1.23") {
		t.Error("error")
	}
}

func TestDestAddr(t *testing.T) {
	dest := "www.baidu.com"
	destAddrs, err := DestAddr(dest)
	if err != nil {
		t.Error(err)
	}

	if len(destAddrs) != 4 {
		t.Error("error")
	}
	t.Log(destAddrs[0], destAddrs[1], destAddrs[2], destAddrs[3])
}

func TestTime2Float(t *testing.T) {
	ff := Time2Float(time.Duration(500) * time.Millisecond)
	t.Log(ff)
	if ff != float32(500) {
		t.Error("error")
	}

}

func TestDestAddrs(t *testing.T) {
	addrs, err := DestAddrs("www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(len(addrs))
	if len(addrs) < 1 {
		t.Error("no ip is error")
	}
	for _, addr := range addrs {
		t.Log(addr)
		if len(addr) != 4 {
			t.Error("ipv4")
		}
	}
}
