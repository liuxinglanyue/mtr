package main

import (
	"fmt"
	"github.com/liuxinglanyue/mtr"
)

func main() {
	fmt.Println("hello")
	fmt.Println(mtr.DEFAULT_RETRIES)
	mtr.LocalAddr()
	destAddrs, _ := mtr.DestAddr("www.baidu.com")

	for _, destAddr := range destAddrs {
		fmt.Println(destAddr)
	}

	//
	c := make(chan mtr.TracerouteHop, 0)
	go func() {
		for {
			hop, ok := <-c
			if !ok {
				fmt.Println()
				return
			}
			fmt.Println(hop.TTL, hop.Address, hop.AvgTime, hop.BestTime, hop.Loss)
		}
	}()
	options := mtr.TracerouteOptions{}
	_, err := mtr.Mtr(destAddrs, &options, c)
	if err != nil {
		fmt.Println(err)
	}
	//

	mm, err := mtr.T("www.baidu.com", true, 0, 0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mm)

	info, err := mtr.T("www.baidu.com", false, 0, 0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}
