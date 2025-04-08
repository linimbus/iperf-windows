package iperf3

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var clientWindow *walk.MainWindow
var clientStatusBar, clientFlowBar *walk.StatusBarItem
var clientAddress *walk.LineEdit
var clientProtocol, clientBandwidthUnit, clientListen *walk.ComboBox
var clientNumberList []*walk.NumberEdit
var clientCheckBoxList []*walk.CheckBox
var clientActive *walk.PushButton
var clientFolder *walk.LineEdit
var clientFolderBut *walk.PushButton
var clientReverseMode *walk.CheckBox
var clientBidirectionalMode *walk.CheckBox
var clientInstance *IperfServer
var clientShutdown bool
var clientRunning bool

func init() {
	clientNumberList = make([]*walk.NumberEdit, 0)
	clientCheckBoxList = make([]*walk.CheckBox, 0)
}

func MakeClientCheckBox(name, tips string, cfg *bool, form walk.Form) CheckBox {
	var box *walk.CheckBox
	return CheckBox{
		AssignTo:      &box,
		Text:          name,
		ToolTipText:   tips,
		Checked:       *cfg,
		StretchFactor: 2,
		OnCheckedChanged: func() {
			*cfg = box.Checked()
			err := configSyncToFile()
			if err != nil {
				ErrorBoxAction(form, err.Error())
			}
		},
		OnBoundsChanged: func() {
			clientCheckBoxList = append(clientCheckBoxList, box)
		},
	}
}

func MakeNumberEdit(max, min int, tips string, cfg *int, form walk.Form) Composite {
	var number *walk.NumberEdit
	return Composite{
		Layout:        HBox{MarginsZero: true},
		StretchFactor: 2,
		Children: []Widget{
			NumberEdit{
				AssignTo:    &number,
				Value:       float64(*cfg),
				ToolTipText: fmt.Sprintf("%d~%d", min, max),
				MaxValue:    float64(max),
				MinValue:    float64(min),
				OnValueChanged: func() {
					*cfg = int(number.Value())
					err := configSyncToFile()
					if err != nil {
						ErrorBoxAction(form, err.Error())
					}
				},
				OnBoundsChanged: func() {
					clientNumberList = append(clientNumberList, number)
				},
			},
			Label{
				Text: tips,
			},
		},
	}
}

func ClientEnable(flag bool) {
	clientAddress.SetEnabled(flag)
	clientProtocol.SetEnabled(flag)
	clientBandwidthUnit.SetEnabled(flag)
	clientListen.SetEnabled(flag)

	for _, but := range clientNumberList {
		but.SetEnabled(flag)
	}

	for _, but := range clientCheckBoxList {
		but.SetEnabled(flag)
	}

	if flag {
		clientActive.SetImage(ICON_Start)
		clientActive.SetToolTipText("Start IPerf3 Client")
		clientActive.SetText("Start")
	} else {
		clientActive.SetImage(ICON_Stop)
		clientActive.SetToolTipText("Stop IPerf3 Client")
		clientActive.SetText("Stop")
	}

	clientFolder.SetEnabled(flag)
	clientFolderBut.SetEnabled(flag)
}

func ClientSwitch() {
	clientActive.SetEnabled(false)
	if clientRunning {
		ClientShutdown()
	} else {
		go ClientActive(configCache)
	}
	time.Sleep(time.Millisecond * 200)
	clientActive.SetEnabled(true)
}

func ClientShutdown() {
	clientShutdown = true
	if clientInstance != nil {
		clientInstance.Shutdown()
	}
	for {
		if !clientRunning {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func ClientActive(config Config) {
	var err error

	defer ClientFlowUpdate("")
	defer ClientEnable(true)

	logs.Info("client active startup")

	repeatCount := config.ClientRepeatCount
	repeatInterval := config.ClientRepeatInterval

	ClientEnable(false)

	clientRunning = true
	for i := 0; i < repeatCount; i++ {
		ClientFlowUpdate(fmt.Sprintf("Repeat Times: %d/%d", i+1, repeatCount))

		if clientShutdown {
			break
		}

		clientInstance, err = ClientStartup(i)
		if err != nil {
			ErrorBoxAction(clientWindow, err.Error())
			break
		}

		for {
			if !clientInstance.running {
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		clientInstance = nil

		if i+1 == repeatCount || clientShutdown {
			break
		}

		time.Sleep(time.Second * time.Duration(repeatInterval))
	}
	clientRunning = false
	clientShutdown = false

	logs.Info("client active stop")

	time.Sleep(time.Millisecond * 200)
}

func ClientFlowUpdate(value string) {
	if clientFlowBar != nil {
		clientFlowBar.SetText(value)
	}
}

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
	if clientInstance != nil {
		clientInstance.Shutdown()
		clientInstance = nil
	}
}

func ClientWindows() error {
	CapSignal(ClientClose)

	interfaces := InterfaceOptions()

	cnt, err := MainWindow{
		Title:    "IPerf3 Client " + VersionGet(),
		Icon:     ICON_Main,
		AssignTo: &clientWindow,
		MinSize:  Size{Width: 600, Height: 250},
		Size:     Size{Width: 600, Height: 250},
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
				Text: "Sponsor",
				OnTriggered: func() {
					AboutAction(clientWindow)
				},
			},
		},
		StatusBarItems: []StatusBarItem{
			{
				AssignTo: &clientStatusBar,
				Text:     cpuInfo(),
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
				Layout: Grid{Columns: 4, Spacing: 6},
				Children: []Widget{
					Label{
						Text:          "Server Address: ",
						ToolTipText:   "Set server ip address to connect",
						StretchFactor: 1,
					},
					LineEdit{
						AssignTo:      &clientAddress,
						Text:          configCache.ClientAddress,
						StretchFactor: 2,
						OnEditingFinished: func() {
							configCache.ClientAddress = clientAddress.Text()
							err := configSyncToFile()
							if err != nil {
								ErrorBoxAction(clientWindow, err.Error())
							}
						},
					},

					Label{
						Text:          "Listen Address: ",
						ToolTipText:   "Set client ip address to bind to",
						StretchFactor: 1,
					},
					ComboBox{
						MaxSize:       Size{Width: 200},
						StretchFactor: 2,
						AssignTo:      &clientListen,
						CurrentIndex:  InterfaceIndex(configCache.ClientListen, interfaces),
						Model:         interfaces,
						OnCurrentIndexChanged: func() {
							configCache.ClientListen = clientListen.Text()
							err := configSyncToFile()
							if err != nil {
								ErrorBoxAction(clientWindow, err.Error())
							}
						},
					},

					Label{
						Text:          "Port: ",
						ToolTipText:   "Set server port to connect",
						StretchFactor: 1,
					},
					MakeNumberEdit(65535, 1, "", &configCache.ClientPort, clientWindow),

					Label{
						Text:          "Protocol: ",
						ToolTipText:   "Using UDP or TCP protocol",
						StretchFactor: 1,
					},
					ComboBox{
						AssignTo:      &clientProtocol,
						StretchFactor: 2,
						CurrentIndex: func() int {
							if configCache.ClientProtocol == "tcp" {
								return 0
							}
							return 1
						}(),
						Model: []string{"tcp", "udp"},
						OnCurrentIndexChanged: func() {
							configCache.ClientProtocol = clientProtocol.Text()
							err := configSyncToFile()
							if err != nil {
								ErrorBoxAction(clientWindow, err.Error())
							}
						},
					},

					Label{
						Text:          "Run Time: ",
						ToolTipText:   "Set the time in seconds to transmit for (default 10 secs)",
						StretchFactor: 1,
					},
					MakeNumberEdit(360000, 0, "Seconds", &configCache.ClientRunTime, clientWindow),

					Label{
						Text:          "Streams: ",
						ToolTipText:   "Set the number of parallel client streams to run",
						StretchFactor: 1,
					},
					MakeNumberEdit(1024, 1, "", &configCache.ClientStreams, clientWindow),

					Label{
						Text:          "Omit Time: ",
						ToolTipText:   "Set perform pre-test for N seconds and omit the pre-test statistics",
						StretchFactor: 1,
					},
					MakeNumberEdit(120, 0, "Seconds", &configCache.ClientOmitSec, clientWindow),

					Label{
						Text:          "Packet Length: ",
						ToolTipText:   "Set the packet length of buffer to read or write (0 as default value 128 KB for TCP, dynamic or 1460 for UDP)",
						StretchFactor: 1,
					},
					MakeNumberEdit(65535, 0, "", &configCache.ClientPayload, clientWindow),

					Label{
						Text:          "IP Dscp: ",
						ToolTipText:   "Set the IP dscp value, either 0-63 or symbolic. Numeric values can be specified in decimal,",
						StretchFactor: 1,
					},
					MakeNumberEdit(63, 0, "", &configCache.ClientDscp, clientWindow),

					Label{
						Text:          "Bandwidth: ",
						ToolTipText:   "Target bitrate in bits/sec (0 for unlimited) (default 1 Mbit/sec for UDP, unlimited for TCP)",
						StretchFactor: 1,
					},
					Composite{
						Layout:        HBox{MarginsZero: true},
						StretchFactor: 2,
						Children: []Widget{
							MakeNumberEdit(999999999, 0, "", &configCache.ClientBandwidth, clientWindow),
							ComboBox{
								AssignTo: &clientBandwidthUnit,
								CurrentIndex: func() int {
									switch configCache.ClientBandwidthUnit {
									case "KB":
										return 0
									case "MB":
										return 1
									case "GB":
										return 2
									default:
										return 0
									}
								}(),
								Model: []string{
									"KB", "MB", "GB",
								},
								OnCurrentIndexChanged: func() {
									configCache.ClientBandwidthUnit = clientBandwidthUnit.Text()
									err := configSyncToFile()
									if err != nil {
										ErrorBoxAction(clientWindow, err.Error())
									}
								},
							},
						},
					},

					Label{
						Text:          "Repeat Count: ",
						StretchFactor: 1,
					},
					MakeNumberEdit(9999999, 1, "", &configCache.ClientRepeatCount, clientWindow),

					Label{
						Text:          "Repeat Interval: ",
						StretchFactor: 1,
					},
					MakeNumberEdit(9999999, 0, "Seconds", &configCache.ClientRepeatInterval, clientWindow),

					Composite{
						Layout:     Grid{Columns: 6, Spacing: 6},
						ColumnSpan: 4,
						Children: []Widget{
							MakeClientCheckBox("No Delay", "Set TCP/SCTP no delay, disabling Nagle's Algorithm", &configCache.ClientNoDelay, clientWindow),
							MakeClientCheckBox("Json Format", "Report in JSON format", &configCache.ClientJsonFormat, clientWindow),
							MakeClientCheckBox("Zero Copy", "Use a 'zero copy' method of sending data", &configCache.ClientZeroCopy, clientWindow),
							MakeClientCheckBox("Dont Fragment", "Set IPv4 Don't Fragment flag", &configCache.ClientDontFragment, clientWindow),

							CheckBox{
								AssignTo:      &clientReverseMode,
								Text:          "Reverse Mode",
								ToolTipText:   "Run in reverse mode, server sends, client receives",
								Checked:       configCache.ClientReverseMode,
								StretchFactor: 2,
								OnCheckedChanged: func() {
									configCache.ClientReverseMode = clientReverseMode.Checked()

									if configCache.ClientReverseMode {
										clientBidirectionalMode.SetChecked(false)
									}

									err := configSyncToFile()
									if err != nil {
										ErrorBoxAction(clientWindow, err.Error())
									}
								},
								OnBoundsChanged: func() {
									clientCheckBoxList = append(clientCheckBoxList, clientReverseMode)
								},
							},
							CheckBox{
								AssignTo:      &clientBidirectionalMode,
								Text:          "Bidirectional Mode",
								ToolTipText:   "Run in bidirectional mode, client and server send and receive",
								Checked:       configCache.ClientBidirectionalMode,
								StretchFactor: 2,
								OnCheckedChanged: func() {
									configCache.ClientBidirectionalMode = clientBidirectionalMode.Checked()

									if configCache.ClientBidirectionalMode {
										clientReverseMode.SetChecked(false)
									}

									err := configSyncToFile()
									if err != nil {
										ErrorBoxAction(clientWindow, err.Error())
									}
								},
								OnBoundsChanged: func() {
									clientCheckBoxList = append(clientCheckBoxList, clientBidirectionalMode)
								},
							},
						},
					},

					Label{
						Text:          "Report Output: ",
						StretchFactor: 1,
					},
					LineEdit{
						ColumnSpan:    2,
						AssignTo:      &clientFolder,
						Text:          configCache.ClientLog,
						StretchFactor: 2,
						OnEditingFinished: func() {
							dir := clientFolder.Text()
							if dir != "" {
								stat, err := os.Stat(dir)
								if err != nil {
									ErrorBoxAction(clientWindow, "The client folder is not exist")
									clientFolder.SetText("")
									dir = ""
								} else if !stat.IsDir() {
									ErrorBoxAction(clientWindow, "The client folder is not directory")
									clientFolder.SetText("")
									dir = ""
								}
							}
							configCache.ClientLog = dir
							err := configSyncToFile()
							if err != nil {
								ErrorBoxAction(clientWindow, err.Error())
							}
						},
					},
					Composite{
						Layout:        HBox{MarginsZero: true},
						StretchFactor: 2,
						Children: []Widget{
							PushButton{
								AssignTo: &clientFolderBut,
								Text:     " ... ",
								OnClicked: func() {
									dlgDir := new(walk.FileDialog)
									dlgDir.FilePath = configCache.ClientLog
									dlgDir.Flags = win.OFN_EXPLORER
									dlgDir.Title = "Please select a folder as output file directory"

									exist, err := dlgDir.ShowBrowseFolder(clientWindow)
									if err != nil {
										logs.Error(err.Error())
										return
									}
									if exist {
										logs.Info("select %s as output file directory", dlgDir.FilePath)

										clientFolder.SetText(dlgDir.FilePath)
										configCache.ClientLog = dlgDir.FilePath
										err := configSyncToFile()
										if err != nil {
											ErrorBoxAction(clientWindow, err.Error())
										}
									}
								},
							},
							PushButton{
								Text: " Open Folder ",
								OnClicked: func() {
									OpenBrowserWeb(configCache.ClientLog)
								},
							},
						},
					},
				},
			},

			Composite{
				Layout: HBox{Margins: Margins{Top: 0, Bottom: 0, Left: 10, Right: 10}},
				Children: []Widget{
					PushButton{
						AssignTo: &clientActive,
						Image:    ICON_Start,
						Text:     "Start",
						OnClicked: func() {
							ClientSwitch()
						},
					},
				},
			},
		},
	}.Run()

	if err != nil {
		logs.Error("client windows startup failed, %s", err.Error())
		return err
	}

	if err := recover(); err != nil {
		logs.Error("main panic, %v", err)
	}

	CloseWindows()

	logs.Info("client windows exit %d", cnt)
	return nil
}
