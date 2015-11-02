package mtr

import (
	"bytes"
	"fmt"
	"net"
)

func T(host string) (result string, err error) {
	options := TracerouteOptions{}

	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return "ip resolve error", err
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("traceroute to %v (%v), %v hops max, %v byte packets\n", host, ipAddr, options.MaxHops(), options.PacketSize()))

	out, err := Traceroute(ipAddr.String(), &options)
	if err == nil {
		if len(out.Hops) == 0 {
			buffer.WriteString("TestTraceroute failed. Expected at least one hop\n")
			return buffer.String(), nil
		}
	} else {
		buffer.WriteString(fmt.Sprintf("TestTraceroute failed due to an error: %v\n", err))
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("%-3v %-16v  %10v  %10v  %10v  %10v  %10v%c\n", "", "HOST", "Avg", "Best", "Wrst", "Last", "Loss", '%'))

	lastTTL := 1
	for _, hop := range out.Hops {
		for j := (lastTTL + 1); j < hop.TTL; j++ {
			buffer.WriteString(fmt.Sprintf("%-3d %-16v  %10.2f  %10.2f  %10.2f  %10.2f  %10.1f%c\n", j, "???", float32(0), float32(0), float32(0), float32(0), float32(100), '%'))
		}
		lastTTL = hop.TTL
		if hop.Success {
			buffer.WriteString(fmt.Sprintf("%-3d %-16v  %10.2f  %10.2f  %10.2f  %10.2f  %10.1f%c\n", hop.TTL, hop.Address, Time2Float(hop.AvgTime), Time2Float(hop.BestTime), Time2Float(hop.WrstTime), Time2Float(hop.LastTime), hop.Loss, '%'))
		} else {
			buffer.WriteString(fmt.Sprintf("%-3d *\n", hop.TTL))
		}
	}

	return buffer.String(), nil
}
