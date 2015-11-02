package mtr

import (
	"net"
	"os"
	"testing"
	"time"
)

const default_local_addr = "0.0.0.0"
const default_dest_addr = "180.97.33.107"

func TestSendIcmp(t *testing.T) {
	ipAddr := net.IPAddr{IP: net.ParseIP(default_dest_addr)}
	pid := os.Getpid() & 0xffff
	ttl := 64
	timeout := time.Duration(500) * time.Millisecond
	hop, err := Icmp(default_local_addr, &ipAddr, ttl, pid, timeout)
	if err != nil {
		t.Error(err)
	}
	t.Log(hop.Success, hop.Addr, hop.Elapsed)
	if !hop.Success {
		t.Error("error, not succ")
	}
	if hop.Addr != default_dest_addr {
		t.Error("error, addr")
	}
}

func TestSendIcmpLoop(t *testing.T) {
	ipAddr := net.IPAddr{IP: net.ParseIP(default_dest_addr)}
	pid := os.Getpid() & 0xffff

	timeout := time.Duration(500) * time.Millisecond

	succ := false

	for ttl := 1; ttl < DEFAULT_MAX_HOPS; ttl++ {
		hop, err := Icmp(default_local_addr, &ipAddr, ttl, pid, timeout)
		if err != nil {
			t.Log("timeout", err)
		}
		t.Log(ttl, hop.Success, hop.Addr, hop.Elapsed)

		if hop.Addr == default_dest_addr {
			succ = true
			break
		}

	}
	if !succ {
		t.Error("error, not succ")
	}

}
