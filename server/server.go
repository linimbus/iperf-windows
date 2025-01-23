package main

import (
	iperf "github.com/linimbus/iperf-windows/iperf3"

	"github.com/astaxie/beego/logs"
)

func main() {
	err := iperf.BoxInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = iperf.FileInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = iperf.LogInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = iperf.IconInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = iperf.ConfigInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	iperf.ServerWindows()
}
