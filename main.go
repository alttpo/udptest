package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var (
	doListen bool
	doSend   bool
	doDump   bool
	nWrite   uint
)

func main() {
	flag.BoolVar(&doListen, "l", false, "listen")
	flag.BoolVar(&doDump, "d", false, "dump hex received")
	flag.BoolVar(&doSend, "s", false, "send")
	flag.UintVar(&nWrite, "n", 65536-12-20, "bytes to send")
	flag.Parse()

	if !doListen && !doSend {
		fmt.Println("use either -l (listen) or -s (send) mode. try -help")
		return
	}

	if doListen {
		listen()
	}
	if doSend {
		send()
	}
}

func listen() {
	var err error
	var uAddr net.UDPAddr
	uAddr.Port = 8193

	var udp *net.UDPConn
	udp, err = net.ListenUDP("udp", &uAddr)
	if err != nil {
		panic(err)
	}

	var b [65536]byte
	for {
		var addr *net.UDPAddr
		var n int
		n, addr, err = udp.ReadFromUDP(b[:])
		if err != nil {
			panic(err)
		}
		_ = addr

		if doDump {
			fmt.Print(hex.Dump(b[:n]))
		} else {
			fmt.Printf("%d\n", n)
		}
	}
}

func send() {
	var err error
	var uAddr net.UDPAddr
	uAddr.Port = 8193

	var udp *net.UDPConn
	udp, err = net.DialUDP("udp", nil, &uAddr)
	if err != nil {
		panic(err)
	}

	// fill b with 65536 bytes
	var b [65536]byte
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(b[:])

	for {
		var n int
		//n, err = udp.Write(b[:65536-20-12])
		n, err = udp.Write(b[:nWrite])
		if err != nil {
			panic(err)
		}

		_ = n
	}
}
