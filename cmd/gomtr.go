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
