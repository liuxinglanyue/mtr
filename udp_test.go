package mtr

import (
	"syscall"
	"testing"
)

func TestSendUdp(t *testing.T) {
	ttl := 3
	local := [4]byte{0, 0, 0, 0}
	dest := [4]byte{180, 97, 30, 107}
	tv := syscall.NsecToTimeval(1000 * 1000 * 500)
	var p = make([]byte, DEFAULT_PACKET_SIZE)

	hop, err := Udp(local, dest, ttl, DEFAULT_PORT, tv, p)

	if err != nil {
		t.Error(err)
	}

	t.Log(hop.Success, hop.Addr, hop.Elapsed)

	if !hop.Success {
		t.Error("error, not succ")
	}

}
