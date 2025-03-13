package main

import (
	iperf "github.com/linimbus/iperf-windows/iperf3"
)

func main() {
	NAME := "client"
	iperf.FileInit(NAME)
	iperf.LogInit(NAME)
	iperf.IconInit()
	iperf.ConfigInit(NAME)
	iperf.ClientWindows()
}
