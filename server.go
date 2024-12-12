package main

import (
	"os"
	"sync"
	"time"

	"github.com/BGrewell/go-iperf"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var serverWindow *walk.MainWindow
var serverInstance *iperf.Server
var serverMutex sync.Mutex
var serverActive, serverFolderBut *walk.PushButton
var serverStatusBar, serverFlowBar *walk.StatusBarItem
var serverPort, serverInterval *walk.NumberEdit
var serverFolder *walk.LineEdit

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

func ServerStatusUpdate(value string) {
	if serverStatusBar != nil {
		serverStatusBar.SetText(value)
	}
}

func ServerRunning() bool {
	return serverInstance != nil
}

func ServerStart() error {
	var err error
	serverInstance, err = ServerStartup()
	if err != nil {
		logs.Warning("iperf server startup failed, %s", err.Error())
		return err
	}
	return nil
}

func ServerStatus(flag bool) {
	serverPort.SetEnabled(!flag)
	serverInterval.SetEnabled(!flag)
}

func ServerShutdown() error {
	if serverInstance != nil {
		serverInstance.Stop()
		serverInstance = nil
	}
	return nil
}

func ServerSwitch() {
	serverMutex.Lock()
	defer serverMutex.Unlock()

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
		serverActive.SetToolTipText("Stop IPerf3 Server")
	} else {
		serverActive.SetImage(ICON_Start)
		serverActive.SetToolTipText("Startup IPerf3 Server")
	}
	ServerStatus(ServerRunning())

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
	ServerShutdown()
}

func ServerWindows() error {
	cnt, err := MainWindow{
		Title:    "IPerf3 Server " + VersionGet(),
		Icon:     ICON_Main,
		AssignTo: &serverWindow,
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
				AssignTo: &serverStatusBar,
				Icon:     ICON_Status,
				Width:    160,
			},
			{
				AssignTo: &serverFlowBar,
				Icon:     ICON_Flow,
				Width:    160,
			},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Listen Port: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							NumberEdit{
								AssignTo: &serverPort,
								Value:    float64(ConfigGet().ServerPort),
								MaxValue: 65535,
								MinValue: 1,
								OnValueChanged: func() {
									err := ServerPortSave(int(serverPort.Value()))
									if err != nil {
										ErrorBoxAction(serverWindow, err.Error())
									}
								},
							},
							Label{
								Text: " 1~65535",
							},
						},
					},
					Label{
						Text: "Statistics Output: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							LineEdit{
								AssignTo: &serverFolder,
								Text:     ConfigGet().ServerLog,
								OnEditingFinished: func() {
									dir := serverFolder.Text()
									if dir != "" {
										stat, err := os.Stat(dir)
										if err != nil {
											ErrorBoxAction(serverWindow, "The server folder is not exist")
											serverFolder.SetText("")
											ServerDirSave("")
											return
										}
										if !stat.IsDir() {
											ErrorBoxAction(serverWindow, "The server folder is not directory")
											serverFolder.SetText("")
											ServerDirSave("")
											return
										}
										return
									}
									ServerDirSave(dir)
								},
							},
							PushButton{
								AssignTo: &serverFolderBut,
								MaxSize:  Size{Width: 30},
								Text:     " ... ",
								OnClicked: func() {
									dlgDir := new(walk.FileDialog)
									dlgDir.FilePath = ConfigGet().ServerLog
									dlgDir.Flags = win.OFN_EXPLORER
									dlgDir.Title = "Please select a folder as output file directory"

									exist, err := dlgDir.ShowBrowseFolder(serverWindow)
									if err != nil {
										logs.Error(err.Error())
										return
									}
									if exist {
										logs.Info("select %s as output file directory", dlgDir.FilePath)
										serverFolder.SetText(dlgDir.FilePath)
										ServerDirSave(dlgDir.FilePath)
									}
								},
							},
						},
					},
					Label{
						Text: "Statistics Interval: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							NumberEdit{
								AssignTo:    &serverInterval,
								Value:       float64(ConfigGet().ServerInterval),
								ToolTipText: "1~60",
								MaxValue:    60,
								MinValue:    1,
								OnValueChanged: func() {
									err := ServerIntervalSave(int(serverInterval.Value()))
									if err != nil {
										ErrorBoxAction(serverWindow, err.Error())
									}
								},
							},
							Label{
								Text: " Seconds",
							},
						},
					},
					Label{
						Text: "Service Options: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							CheckBox{
								Text: "Auto Startup",
							},
							CheckBox{
								Text: "Json Format",
							},
						},
					},
					HSpacer{},
					PushButton{
						AssignTo:    &serverActive,
						Image:       ICON_Start,
						MinSize:     Size{Width: 200},
						Text:        " ",
						ToolTipText: "Startup IPerf3 Server",
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
