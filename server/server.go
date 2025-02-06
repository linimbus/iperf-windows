package main

import (
	iperf "github.com/linimbus/iperf-windows/iperf3"
)

func main() {
	NAME := "server"
	iperf.FileInit()
	iperf.LogInit(NAME)
	iperf.IconInit()
	iperf.ConfigInit(NAME)
	iperf.ServerWindows()
}
