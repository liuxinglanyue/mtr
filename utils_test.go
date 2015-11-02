package mtr

import (
	"testing"
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
