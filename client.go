package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var clientWindow *walk.MainWindow

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
	var statusBar *walk.StatusBarItem

	cnt, err := MainWindow{
		Title:    "IPerf3 Client " + VersionGet(),
		Icon:     ICON_Main,
		AssignTo: &clientWindow,
		MinSize:  Size{Width: 200, Height: 300},
		Size:     Size{Width: 200, Height: 300},
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
				AssignTo: &statusBar,
				Text:     "",
				Icon:     ICON_Status,
				Width:    300,
				OnClicked: func() {
				},
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
