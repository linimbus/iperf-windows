package iperf3

import (
	"os"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var serverWindow *walk.MainWindow
var serverInstance *IperfServer
var serverMutex sync.Mutex
var serverActive, serverFolderBut *walk.PushButton
var serverStatusBar, serverFlowBar *walk.StatusBarItem
var serverPort, serverInterval *walk.NumberEdit
var serverFolder *walk.LineEdit

var serverCheckBoxList []*walk.CheckBox

func MakeCheckBox(name string, cfg *bool, form walk.Form) CheckBox {
	var box *walk.CheckBox
	return CheckBox{
		AssignTo: &box,
		Text:     name,
		Checked:  *cfg,
		OnCheckedChanged: func() {
			*cfg = box.Checked()
			err := configSyncToFile()
			if err != nil {
				ErrorBoxAction(form, err.Error())
			}
		},
		OnBoundsChanged: func() {
			serverCheckBoxList = append(serverCheckBoxList, box)
		},
	}
}

func init() {
	serverCheckBoxList = make([]*walk.CheckBox, 0)
	go func() {
		for {
			if serverWindow != nil && serverWindow.Visible() {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		NotifyInit()

		if configCache.ServerAutoHide {
			serverWindow.SetVisible(false)
		}

		if configCache.ServerAutoStartup && !ServerRunning() {
			ServerSwitch()
		}

		for {
			time.Sleep(time.Millisecond * 100)

			serverMutex.Lock()
			if serverInstance != nil && !serverInstance.running {
				ServerStatus(ServerRunning())
			}
			serverMutex.Unlock()
		}
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
	serverFolderBut.SetEnabled(!flag)
	serverFolder.SetEnabled(!flag)
	for _, box := range serverCheckBoxList {
		box.SetEnabled(!flag)
	}
	if flag {
		serverActive.SetImage(ICON_Stop)
		serverActive.SetToolTipText("Stop IPerf3 Server")
		serverActive.SetText("Stop")
	} else {
		serverActive.SetImage(ICON_Start)
		serverActive.SetToolTipText("Start IPerf3 Server")
		serverActive.SetText("Start")
	}
}

func ServerShutdown() error {
	if serverInstance != nil {
		serverInstance.Shutdown()
		serverInstance = nil
	}
	return nil
}

func ServerSwitch() {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	defer serverActive.SetEnabled(true)

	var err error
	if ServerRunning() {
		err = ServerShutdown()
	} else {
		err = ServerStart()
	}

	if err != nil {
		ErrorBoxAction(serverWindow, err.Error())
	}

	ServerStatus(ServerRunning())
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
	CapSignal(CloseWindows)

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
				Text: "Hide Window",
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
				Text:     cpuInfo(),
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
						Text: "Report Output: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							LineEdit{
								AssignTo: &serverFolder,
								Text:     configCache.ServerLog,
								OnEditingFinished: func() {
									dir := serverFolder.Text()
									if dir != "" {
										stat, err := os.Stat(dir)
										if err != nil {
											ErrorBoxAction(serverWindow, "The server folder is not exist")
											serverFolder.SetText("")
											dir = ""
										} else if !stat.IsDir() {
											ErrorBoxAction(serverWindow, "The server folder is not directory")
											serverFolder.SetText("")
											dir = ""
										}
									}
									configCache.ServerLog = dir
									err := configSyncToFile()
									if err != nil {
										ErrorBoxAction(serverWindow, err.Error())
									}
								},
							},
							PushButton{
								AssignTo: &serverFolderBut,
								Text:     "...",
								OnClicked: func() {
									dlgDir := new(walk.FileDialog)
									dlgDir.FilePath = configCache.ServerLog
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
										configCache.ServerLog = dlgDir.FilePath
										err := configSyncToFile()
										if err != nil {
											ErrorBoxAction(serverWindow, err.Error())
										}
									}
								},
							},
							PushButton{
								Text: "Open",
								OnClicked: func() {
									OpenBrowserWeb(configCache.ServerLog)
								},
							},
						},
					},
					Label{
						Text: "Report Interval: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							NumberEdit{
								AssignTo:    &serverInterval,
								Value:       float64(configCache.ServerInterval),
								ToolTipText: "1~60",
								MaxValue:    60,
								MinValue:    1,
								OnValueChanged: func() {
									configCache.ServerInterval = int(serverInterval.Value())
									err := configSyncToFile()
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
						Text: "Service Port: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							NumberEdit{
								AssignTo: &serverPort,
								Value:    float64(configCache.ServerPort),
								MaxValue: 65535,
								MinValue: 1,
								OnValueChanged: func() {
									configCache.ServerPort = int(serverPort.Value())
									err := configSyncToFile()
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
						Text: "Service Options: ",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							MakeCheckBox("Auto Startup", &configCache.ServerAutoStartup, serverWindow),
							MakeCheckBox("Auto Hide", &configCache.ServerAutoHide, serverWindow),
							MakeCheckBox("Json Format", &configCache.ServerJsonFormat, serverWindow),
						},
					},
				},
			},
			Composite{
				Layout: VBox{Margins: Margins{Top: 0, Bottom: 0, Left: 10, Right: 10}},
				Children: []Widget{
					PushButton{
						AssignTo: &serverActive,
						Image:    ICON_Start,
						MinSize:  Size{Width: 200},
						Text:     "Start",
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

	CloseWindows()

	return nil
}
