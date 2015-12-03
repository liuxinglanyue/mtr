package mtr

import (
	"net"
	"os"
	"time"
)

func Mtr(dest [4]byte, options *TracerouteOptions, c ...chan TracerouteHop) (result TracerouteResult, err error) {
	result.Hops = []TracerouteHop{}
	result.DestAddress = dest
	destAddr := AddressString(dest)
	localAddr := "0.0.0.0"
	ipAddr := net.IPAddr{IP: net.ParseIP(destAddr)}
	pid := os.Getpid() & 0xffff
	timeout := time.Duration(options.TimeoutMs()) * time.Millisecond

	mtrResults := make([]MtrResult, options.MaxHops()+1)

	for snt := 0; snt < options.SntSize(); snt++ {
		retry := 0
		for ttl := 1; ttl < options.MaxHops(); ttl++ {
			snt_plus := snt + 1
			hop := TracerouteHop{TTL: ttl, Snt: snt_plus}
			if mtrResults[ttl].TTL == 0 {
				mtrResults[ttl] = MtrResult{TTL: ttl, Host: "???", SuccSum: 0, Success: false, LastTime: time.Duration(0), AllTime: time.Duration(0), BestTime: time.Duration(0), WrstTime: time.Duration(0), AvgTime: time.Duration(0)}
			}
			hopReturn, err := Icmp(localAddr, &ipAddr, ttl, pid, timeout)
			if err != nil || !hopReturn.Success {
				//mtrResults[ttl].Success = false
				hop.Loss = (float32)(snt+1-mtrResults[ttl].SuccSum) / (float32)(snt+1) * 100
				hop.Address = mtrResults[ttl].Host
				hop.AvgTime = mtrResults[ttl].AvgTime
				hop.BestTime = mtrResults[ttl].BestTime
				hop.Host = mtrResults[ttl].Host
				hop.LastTime = mtrResults[ttl].LastTime
				hop.Success = false
				hop.WrstTime = mtrResults[ttl].WrstTime
				notifyMtr(hop, c)
				retry++
				if retry >= options.Retries() {
					break
				}
				continue
			}
			retry = 0
			mtrResults[ttl].SuccSum = mtrResults[ttl].SuccSum + 1
			mtrResults[ttl].Host = hopReturn.Addr
			mtrResults[ttl].LastTime = hopReturn.Elapsed
			if mtrResults[ttl].WrstTime == time.Duration(0) || hopReturn.Elapsed > mtrResults[ttl].WrstTime {
				mtrResults[ttl].WrstTime = hopReturn.Elapsed
			}
			if mtrResults[ttl].BestTime == time.Duration(0) || hopReturn.Elapsed < mtrResults[ttl].BestTime {
				mtrResults[ttl].BestTime = hopReturn.Elapsed
			}
			mtrResults[ttl].AllTime += hopReturn.Elapsed
			mtrResults[ttl].AvgTime = time.Duration((int64)(mtrResults[ttl].AllTime/time.Microsecond)/(int64)(mtrResults[ttl].SuccSum)) * time.Microsecond

			hop.Loss = (float32)(snt+1-mtrResults[ttl].SuccSum) / (float32)(snt+1) * 100
			hop.Address = mtrResults[ttl].Host
			hop.AvgTime = mtrResults[ttl].AvgTime
			hop.BestTime = mtrResults[ttl].BestTime
			hop.Host = mtrResults[ttl].Host
			hop.LastTime = mtrResults[ttl].LastTime
			hop.Success = true
			hop.WrstTime = mtrResults[ttl].WrstTime
			notifyMtr(hop, c)

			if hop.Address == destAddr {
				break
			}
		}
	}

	retry := 0
	for _, mtrResult := range mtrResults {
		if mtrResult.TTL == 0 {
			retry++
			if retry >= options.Retries() {
				break
			}
			continue
		}
		retry = 0
		hop := TracerouteHop{TTL: mtrResult.TTL}
		hop.Address = mtrResult.Host
		hop.Host = mtrResult.Host
		hop.AvgTime = mtrResult.AvgTime
		hop.BestTime = mtrResult.BestTime
		hop.LastTime = mtrResult.LastTime
		failSum := options.SntSize() - mtrResult.SuccSum
		loss := (float32)(failSum) / (float32)(options.SntSize()) * 100
		hop.Loss = float32(loss)
		hop.WrstTime = mtrResult.WrstTime
		hop.Success = true

		result.Hops = append(result.Hops, hop)

		if hop.Host == destAddr {
			break
		}
	}

	closeNotify(c)

	return result, nil
}

func notifyMtr(hop TracerouteHop, channels []chan TracerouteHop) {
	for _, c := range channels {
		c <- hop
	}
}

func closeNotifyMtr(channels []chan TracerouteHop) {
	for _, c := range channels {
		close(c)
	}
}
