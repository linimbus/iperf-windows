package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var clientWindow *walk.MainWindow
var clientStatusBar, clientFlowBar *walk.StatusBarItem

func ClientStatusUpdate(value string) {
	if clientStatusBar != nil {
		clientStatusBar.SetText(value)
	}
}

func ClientClose() {
	if clientWindow != nil {
		err := clientWindow.Close()
		if err != nil {
			logs.Warning("client close %s", err.Error())
		}
		clientWindow = nil
	}
}

func ClientWindows() error {

	cnt, err := MainWindow{
		Title:    "IPerf3 Client " + VersionGet(),
		Icon:     ICON_Main,
		AssignTo: &clientWindow,
		MinSize:  Size{Width: 400, Height: 250},
		Size:     Size{Width: 400, Height: 250},
		Layout:   VBox{Margins: Margins{Top: 5, Bottom: 5, Left: 5, Right: 5}},
		Font:     Font{Bold: true},
		MenuItems: []MenuItem{
			Action{
				Text: "Runlog",
				OnTriggered: func() {
					OpenBrowserWeb(RunlogDirGet())
				},
			},
			Action{
				Text: "Mini Windows",
				OnTriggered: func() {
					NotifyInit()
					clientWindow.SetVisible(false)
				},
			},
			Action{
				Text: "Sponsor",
				OnTriggered: func() {
					AboutAction(clientWindow)
				},
			},
		},
		StatusBarItems: []StatusBarItem{
			{
				AssignTo: &clientStatusBar,
				Icon:     ICON_Status,
				Width:    160,
			},
			{
				AssignTo: &clientFlowBar,
				Icon:     ICON_Flow,
				Width:    160,
			},
		},
		Children: []Widget{
			Composite{
				Layout:   Grid{Columns: 2},
				Children: []Widget{},
			},
		},
	}.Run()

	if err != nil {
		logs.Error("client windows startup failed, %s", err.Error())
		return err
	}

	logs.Info("client windows exit %d", cnt)

	shutdown <- struct{}{}

	return nil
}
