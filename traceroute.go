package mtr

import (
	"syscall"
	"time"
)

func Traceroute(dest string, options *TracerouteOptions, c ...chan TracerouteHop) (result TracerouteResult, err error) {
	result.Hops = []TracerouteHop{}
	destAddr, err := DestAddr(dest)
	result.DestAddress = destAddr
	socketAddr, err := LocalAddr()
	if err != nil {
		return
	}

	timeoutMs := (int64)(options.TimeoutMs())
	tv := syscall.NsecToTimeval(1000 * 1000 * timeoutMs)
	var p = make([]byte, options.PacketSize())

	retry := 0
	succ := false
	sntCount := options.SntSize()
	retryCount := options.Retries()

	for ttl := 1; ttl < options.MaxHops(); ttl++ {
		if ttl >= 10 && retry > 1 {
			sntCount = options.SntSize() / retry
			retryCount = options.Retries() / 2
			if sntCount < 2 {
				sntCount = 2
			}
		}
		hop := TracerouteHop{TTL: ttl}
		failSum := 0
		allTime := time.Duration(0) * time.Second
		firstReturn := true
		for i := 0; i < sntCount; i++ {
			singleHop, err := Udp(socketAddr, destAddr, ttl, options.Port(), tv, p)
			if err != nil {
				failSum++
				continue
			}
			hop.LastTime = singleHop.Elapsed
			allTime += singleHop.Elapsed

			if firstReturn {
				hop.Address = singleHop.Addr
				hop.BestTime = singleHop.Elapsed
				hop.WrstTime = singleHop.Elapsed
				firstReturn = false
			} else {
				if singleHop.Elapsed > hop.WrstTime {
					hop.WrstTime = singleHop.Elapsed
				} else if singleHop.Elapsed < hop.BestTime {
					hop.BestTime = singleHop.Elapsed
				}
			}

			if singleHop.Addr == AddressString(destAddr) {
				succ = true
			}
		}
		if failSum == sntCount {
			hop.Success = false
			retry++
			if retry >= retryCount {
				closeNotify(c)
				return result, nil
			}
			continue
		}

		retry = 0
		hop.AvgTime = time.Duration((int64)(allTime/time.Microsecond)/(int64)(sntCount-failSum)) * time.Microsecond
		hop.Loss = float32(failSum) / float32(sntCount) * 100
		hop.Success = true
		result.Hops = append(result.Hops, hop)

		notify(hop, c)

		if succ {
			closeNotify(c)
			return result, nil
		}
	}
	closeNotify(c)
	return result, nil
}

func notify(hop TracerouteHop, channels []chan TracerouteHop) {
	for _, c := range channels {
		c <- hop
	}
}

func closeNotify(channels []chan TracerouteHop) {
	for _, c := range channels {
		close(c)
	}
}
