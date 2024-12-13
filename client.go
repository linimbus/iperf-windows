package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var clientWindow *walk.MainWindow
var clientStatusBar, clientFlowBar *walk.StatusBarItem
var clientAddress *walk.LineEdit
var clientProtocol, clientBandwidthUnit *walk.ComboBox
var clientProcessBar *walk.ProgressBar
var clientNumberList []*walk.NumberEdit
var clientCheckBoxList []*walk.CheckBox
var clientActive *walk.PushButton

func init() {
	clientNumberList = make([]*walk.NumberEdit, 0)
	clientCheckBoxList = make([]*walk.CheckBox, 0)
}

func MakeClientCheckBox(name, tips string, cfg *bool, form walk.Form) CheckBox {
	var box *walk.CheckBox
	return CheckBox{
		AssignTo:    &box,
		Text:        name,
		ToolTipText: tips,
		Checked:     *cfg,
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
		Layout: HBox{MarginsZero: true},
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
	clientActive.SetEnabled(flag)
	clientAddress.SetEnabled(flag)
	clientProtocol.SetEnabled(flag)
	clientBandwidthUnit.SetEnabled(flag)

	for _, but := range clientNumberList {
		but.SetEnabled(flag)
	}

	for _, but := range clientCheckBoxList {
		but.SetEnabled(flag)
	}
}

func ClientActive() {
	defer ClientEnable(true)

	c, err := ClientStartup()
	if err != nil {
		ErrorBoxAction(clientWindow, err.Error())
		return
	}

	for i := 0; i < configCache.ClientRunTime+1; i++ {
		time.Sleep(time.Second)
		clientProcessBar.SetValue((i * 100) / (configCache.ClientRunTime + 1))
	}

	for {
		clientProcessBar.SetValue(100)
		if !c.running {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	clientProcessBar.SetValue(0)
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
				Text: "Hide Window",
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
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text:        "Server Address: ",
						ToolTipText: "Set server ip address to connect",
					},
					LineEdit{
						AssignTo: &clientAddress,
						Text:     configCache.ClientAddress,
						OnEditingFinished: func() {
							configCache.ClientAddress = clientAddress.Text()
							err := configSyncToFile()
							if err != nil {
								ErrorBoxAction(clientWindow, err.Error())
							}
						},
					},
					Label{
						Text:        "Port: ",
						ToolTipText: "Set server port to connect",
					},
					MakeNumberEdit(65535, 1, "1~65535", &configCache.ClientPort, clientWindow),
					Label{
						Text:        "Protocol: ",
						ToolTipText: "Using UDP or TCP protocol",
					},
					ComboBox{
						AssignTo: &clientProtocol,
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
						Text:        "Run Time: ",
						ToolTipText: "Set the time in seconds to transmit for (default 10 secs)",
					},
					MakeNumberEdit(3600, 1, "Seconds", &configCache.ClientRunTime, clientWindow),

					Label{
						Text:        "Streams: ",
						ToolTipText: "Set the number of parallel client streams to run",
					},
					MakeNumberEdit(1024, 1, "", &configCache.ClientStreams, clientWindow),

					Label{
						Text:        "Omit Time: ",
						ToolTipText: "Set perform pre-test for N seconds and omit the pre-test statistics",
					},
					MakeNumberEdit(120, 0, "Seconds", &configCache.ClientOmitSec, clientWindow),

					Label{
						Text:        "Packet Length: ",
						ToolTipText: "Set the packet length of buffer to read or write (0 as default value 128 KB for TCP, dynamic or 1460 for UDP)",
					},
					MakeNumberEdit(65535, 0, "", &configCache.ClientPayload, clientWindow),

					MakeClientCheckBox("Dont Fragment", "Set IPv4 Don't Fragment flag", &configCache.ClientDontFragment, clientWindow),
					MakeClientCheckBox("Json Format", "Report in JSON format", &configCache.ClientJsonFormat, clientWindow),

					Label{
						Text:        "IP Dscp: ",
						ToolTipText: "Set the IP dscp value, either 0-63 or symbolic. Numeric values can be specified in decimal,",
					},
					MakeNumberEdit(63, 0, "", &configCache.ClientDscp, clientWindow),

					MakeClientCheckBox("Zero Copy", "Use a 'zero copy' method of sending data", &configCache.ClientZeroCopy, clientWindow),
					MakeClientCheckBox("No Delay", "Set TCP/SCTP no delay, disabling Nagle's Algorithm", &configCache.ClientNoDelay, clientWindow),

					Label{
						Text:        "Bandwidth: ",
						ToolTipText: "Target bitrate in bits/sec (0 for unlimited) (default 1 Mbit/sec for UDP, unlimited for TCP)",
					},
					Composite{
						Layout: HBox{MarginsZero: true},
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
				},
			},
			Composite{
				Layout: HBox{Margins: Margins{Top: 0, Bottom: 0, Left: 12, Right: 12}},
				Children: []Widget{
					ProgressBar{
						AssignTo: &clientProcessBar,
						MaxValue: 100,
						MinValue: 0,
						MinSize:  Size{Height: 2},
						MaxSize:  Size{Height: 2},
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
							ClientEnable(false)

							go ClientActive()
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

	logs.Info("client windows exit %d", cnt)

	shutdown <- struct{}{}

	return nil
}
