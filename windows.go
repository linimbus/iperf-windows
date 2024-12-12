package main

import (
	"time"

	"github.com/astaxie/beego/logs"
)

var shutdown chan struct{}

func init() {
	shutdown = make(chan struct{}, 10)
}

func mainWindows() {
	CapSignal(CloseWindows)

	go ServerWindows()
	time.Sleep(time.Millisecond * 100)

	go ClientWindows()
	time.Sleep(time.Millisecond * 100)

	<-shutdown

	if err := recover(); err != nil {
		logs.Error("main panic, %v", err)
	}

	logs.Info("main windows shutdown")

	CloseWindows()
}

func CloseWindows() {
	ClientClose()
	ServerClose()
	NotifyExit()
	shutdown <- struct{}{}
}
