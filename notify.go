package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

var notify *walk.NotifyIcon

func NotifyExit() {
	if notify != nil {
		err := notify.Dispose()
		if err != nil {
			logs.Warning("notify dispose failed, %s", err.Error())
		}
		notify = nil
	}
}

var lastCheck time.Time

func NotifyInit() {
	if notify != nil {
		return
	}

	var err error

	notify, err = walk.NewNotifyIcon(serverWindow)
	if err != nil {
		logs.Error("new notify icon fail, %s", err.Error())
		return
	}

	err = notify.SetIcon(ICON_Main)
	if err != nil {
		logs.Error("set notify icon fail, %s", err.Error())
		return
	}

	err = notify.SetToolTip("")
	if err != nil {
		logs.Error("set notify tool tip fail, %s", err.Error())
		return
	}

	exitBut := walk.NewAction()
	err = exitBut.SetText("Exit")
	if err != nil {
		logs.Error("notify new action fail, %s", err.Error())
		return
	}

	exitBut.Triggered().Attach(func() {
		CloseWindows()
	})

	serverBut := walk.NewAction()
	err = serverBut.SetText("Show Server Windows")
	if err != nil {
		logs.Error("notify new action fail, %s", err.Error())
		return
	}

	serverBut.Triggered().Attach(func() {
		serverWindow.SetVisible(true)
	})

	clientBut := walk.NewAction()
	err = clientBut.SetText("Show Client Windows")
	if err != nil {
		logs.Error("notify new action fail, %s", err.Error())
		return
	}

	clientBut.Triggered().Attach(func() {
		clientWindow.SetVisible(true)
	})

	if err := notify.ContextMenu().Actions().Add(clientBut); err != nil {
		logs.Error("notify add action fail, %s", err.Error())
		return
	}

	if err := notify.ContextMenu().Actions().Add(serverBut); err != nil {
		logs.Error("notify add action fail, %s", err.Error())
		return
	}

	if err := notify.ContextMenu().Actions().Add(exitBut); err != nil {
		logs.Error("notify add action fail, %s", err.Error())
		return
	}

	notify.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		now := time.Now()
		if now.Sub(lastCheck) < time.Second {
			clientWindow.SetVisible(true)
			serverWindow.SetVisible(true)
		}
		lastCheck = now
	})

	notify.SetVisible(true)
}
