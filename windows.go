package main

import (
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 300
var mainWindowHeight = 200

func init() {
	go func() {
		for {
			if mainWindow != nil && mainWindow.Visible() {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		NotifyAction()
	}()
}

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Action{
			Text: "Runlog",
			OnTriggered: func() {
				OpenBrowserWeb(RunlogDirGet())
			},
		},
		Action{
			Text: "Mini Windows",
			OnTriggered: func() {
				NotifyAction()
			},
		},
		Action{
			Text: "Sponsor",
			OnTriggered: func() {
				AboutAction()
			},
		},
	}
}

var intervalNumber *walk.NumberEdit
var filterInterface, restfulMethod *walk.ComboBox
var connectivityURL, restfulURL, restfulHeaderKey, restfulHeaderValue, outputFolder *walk.LineEdit

func ConsoleWidget() []Widget {
	httpMethodList := []string{
		"PUT", "POST", "GET", "PATCH", "HEAD", "DELETE",
	}
	interfaceList := InterfaceOptions()

	return []Widget{
		Label{
			Text: "Connectivity URL: ",
		},
		LineEdit{
			AssignTo: &connectivityURL,
			Text:     ConfigGet().ConnectivityURL,
			OnEditingFinished: func() {
				ConnectivityURLSave(connectivityURL.Text())
			},
		},
		Label{
			Text: "Interface Monitor: ",
		},
		ComboBox{
			AssignTo: &filterInterface,
			Model:    interfaceList,
			CurrentIndex: func() int {
				filter := ConfigGet().FilterInterface
				for i, name := range interfaceList {
					if name == filter {
						return i
					}
				}
				return 0
			},
			OnCurrentIndexChanged: func() {
				err := FilterInterfaceSave(filterInterface.Text())
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
				}
			},
			OnBoundsChanged: func() {
				addr := ConfigGet().FilterInterface
				for i, item := range interfaceList {
					if addr == item {
						filterInterface.SetCurrentIndex(i)
						return
					}
				}
				filterInterface.SetCurrentIndex(0)
			},
		},
		Label{
			Text: "Output Folder: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &outputFolder,
					Text:     ConfigGet().OutputDirectory,
					OnEditingFinished: func() {
						dir := outputFolder.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "The server folder is not exist")
								outputFolder.SetText("")
								OutputDirectorySave("")
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "The server folder is not directory")
								outputFolder.SetText("")
								OutputDirectorySave("")
								return
							}
						}
						OutputDirectorySave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 30},
					Text:    " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().OutputDirectory
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as output file directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as output file directory", dlgDir.FilePath)
							outputFolder.SetText(dlgDir.FilePath)
							OutputDirectorySave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Label{
			Text: "Restful URL: ",
		},
		LineEdit{
			AssignTo: &restfulURL,
			Text:     ConfigGet().RestfulURL,
			OnEditingFinished: func() {
				RestfulURLSave(restfulURL.Text())
			},
		},
		Label{
			Text: "Restful Method: ",
		},
		ComboBox{
			AssignTo: &restfulMethod,
			Model:    httpMethodList,
			CurrentIndex: func() int {
				filter := ConfigGet().RestfulMethod
				for i, name := range httpMethodList {
					if name == filter {
						return i
					}
				}
				return 0
			},
			OnCurrentIndexChanged: func() {
				err := RestfulMethodSave(restfulMethod.Text())
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
				}
			},
			OnBoundsChanged: func() {
				addr := ConfigGet().RestfulMethod
				for i, item := range httpMethodList {
					if addr == item {
						restfulMethod.SetCurrentIndex(i)
						return
					}
				}
				restfulMethod.SetCurrentIndex(0)
			},
		},
		Label{
			Text: "Restful Header: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &restfulHeaderKey,
					Text: func() string {
						list := strings.Split(ConfigGet().RestfulHeader, ":")
						return list[0]
					}(),
					OnEditingFinished: func() {
						RestfulHeaderSave(restfulHeaderKey.Text(), restfulHeaderValue.Text())
					},
				},
				Label{
					Text: ":",
				},
				LineEdit{
					AssignTo: &restfulHeaderValue,
					Text: func() string {
						list := strings.Split(ConfigGet().RestfulHeader, ":")
						if len(list) == 2 {
							return list[1]
						}
						return ""
					}(),
					OnEditingFinished: func() {
						RestfulHeaderSave(restfulHeaderKey.Text(), restfulHeaderValue.Text())
					},
				},
			},
		},

		Label{
			Text: "Loop Interval: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				NumberEdit{
					AssignTo:    &intervalNumber,
					Value:       float64(ConfigGet().Interval),
					ToolTipText: "5~300",
					MaxValue:    300,
					MinValue:    5,
					OnValueChanged: func() {
						err := IntervalSave(int(intervalNumber.Value()))
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
					},
				},
				Label{
					Text: "Seconds",
				},
			},
		},
	}
}

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:          "Simple IP Monitor " + VersionGet(),
		Icon:           ICON_Main,
		AssignTo:       &mainWindow,
		MinSize:        Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Size:           Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Layout:         VBox{Margins: Margins{Top: 5, Bottom: 5, Left: 5, Right: 5}},
		Font:           Font{Bold: true},
		MenuItems:      MenuBarInit(),
		StatusBarItems: StatusBarInit(),
		Children: []Widget{
			Composite{
				Layout:   Grid{Columns: 2},
				Children: ConsoleWidget(),
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	if err := recover(); err != nil {
		logs.Error(err)
	}

	CloseWindows()
}

func CloseWindows() {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
	NotifyExit()
}
