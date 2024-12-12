package main

import (
	"sync"
	"time"

	"github.com/BGrewell/go-iperf"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var serverWindow *walk.MainWindow
var serverInstance *iperf.Server
var serverMutex sync.Mutex
var serverActive *walk.PushButton

func init() {
	go func() {
		for {
			if serverWindow != nil && serverWindow.Visible() {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		// NotifyAction()
	}()
}

func ServerRunning() bool {
	return serverInstance != nil
}

func ServerStart() error {
	var err error
	serverInstance, err = ServerStartup(5012)
	if err != nil {
		logs.Warning("iperf server startup failed, %s", err.Error())
		return err
	}
	return nil
}

func ServerShutdown() error {
	serverInstance.Stop()
	serverInstance = nil
	return nil
}

func ServerSwitch() {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	time.Sleep(time.Millisecond * 500)

	var err error
	if ServerRunning() {
		err = ServerShutdown()
	} else {
		err = ServerStart()
	}

	if err != nil {
		ErrorBoxAction(serverWindow, err.Error())
	}

	if ServerRunning() {
		serverActive.SetImage(ICON_Stop)
	} else {
		serverActive.SetImage(ICON_Start)
	}
	// ServerStatus(ServerRunning())

	serverActive.SetEnabled(true)
}

func ServerClose() {
	if serverWindow != nil {
		err := serverWindow.Close()
		if err != nil {
			logs.Warning("server close %s", err.Error())
		}
		serverWindow = nil
	}
}

func ServerWindows() error {

	var listenPort *walk.NumberEdit
	var listenAddr *walk.ComboBox
	var statusBar *walk.StatusBarItem

	interfaceList := InterfaceOptions()

	cnt, err := MainWindow{
		Title:    "IPerf3 Server " + VersionGet(),
		Icon:     ICON_Main,
		AssignTo: &serverWindow,
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
					serverWindow.SetVisible(false)
				},
			},
			Action{
				Text: "Sponsor",
				OnTriggered: func() {
					AboutAction(serverWindow)
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
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Listen: ",
					},
					ComboBox{
						AssignTo: &listenAddr,
						CurrentIndex: func() int {
							// addr := ConfigGet().ListenAddr
							// for i, item := range interfaceList {
							// 	if addr == item {
							// 		return i
							// 	}
							// }
							return 0
						},
						Model: interfaceList,
						OnCurrentIndexChanged: func() {
							// err := ListenAddressSave(listenAddr.Text())
							// if err != nil {
							// 	ErrorBoxAction(mainWindow, err.Error())
							// } else {
							// 	BrowseURLUpdate()
							// }
						},
						OnBoundsChanged: func() {
							// addr := ConfigGet().ListenAddr
							// for i, item := range interfaceList {
							// 	if addr == item {
							// 		listenAddr.SetCurrentIndex(i)
							// 		return
							// 	}
							// }
							listenAddr.SetCurrentIndex(0)
						},
					},
					Label{
						Text: "Port: ",
					},
					NumberEdit{
						AssignTo:    &listenPort,
						Value:       float64(8090),
						ToolTipText: "1~65535",
						MaxValue:    65535,
						MinValue:    1,
						OnValueChanged: func() {
							// err := ListenPortSave(int64(listenPort.Value()))
							// if err != nil {
							// 	ErrorBoxAction(mainWindow, err.Error())
							// } else {
							// 	BrowseURLUpdate()
							// }
						},
					},
					HSpacer{},
					PushButton{
						AssignTo:    &serverActive,
						Image:       ICON_Start,
						Text:        " ",
						ToolTipText: "Startup or Stop",
						MinSize:     Size{Height: 32},
						OnClicked: func() {
							serverActive.SetEnabled(false)
							go ServerSwitch()
						},
					},
				},
			},
		},
	}.Run()

	if err != nil {
		logs.Error("server windows startup failed, %s", err.Error())
		return err
	}

	logs.Info("server windows exit %d", cnt)

	shutdown <- struct{}{}

	return nil
}
