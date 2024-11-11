package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func VersionGet() string {
	return "v1.0.0"
}

func main() {

	var mainWin *walk.MainWindow

	var consoleMode *walk.ComboBox
	var consoleAddress *walk.LineEdit
	var consolePort *walk.NumberEdit
	var consoleStreams *walk.NumberEdit
	var consoleLength *walk.NumberEdit
	var consoleTimeout *walk.NumberEdit
	var consoleProtocol *walk.ComboBox
	var consoleVersion *walk.ComboBox
	var buttonStart, buttonStop *walk.PushButton

	_, err := MainWindow{
		Title:    "iperf-windows " + VersionGet(),
		AssignTo: &mainWin,
		Icon:     walk.IconApplication(),
		MinSize:  Size{Width: 350, Height: 350},
		Size:     Size{Width: 350, Height: 350},
		Font:     Font{Family: "Segoe UI", PointSize: 10},
		Layout: VBox{
			Alignment:   AlignHNearVNear,
			MarginsZero: true,
			Margins:     Margins{Left: 10, Top: 10, Right: 10, Bottom: 10},
		},
		MenuItems: []MenuItem{
			Action{
				Text: "Setting",
				OnTriggered: func() {

				},
			},
			Action{
				Text: "Mini Windows",
				OnTriggered: func() {

				},
			},
			Action{
				Text: "About",
				OnTriggered: func() {
					AboutAction(mainWin)
				},
			},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Mode: ",
					},
					ComboBox{
						AssignTo:     &consoleMode,
						CurrentIndex: 0,
						Model:        []string{"Server", "Client"},
						OnCurrentIndexChanged: func() {
							if consoleMode.Text() == "Client" {
								mainWin.SetSize(walk.Size{Width: 350, Height: 500})
							} else {
								mainWin.SetSize(walk.Size{Width: 350, Height: 350})
							}
						},
					},
					Label{
						Text: "Address:",
					},
					LineEdit{
						AssignTo:  &consoleAddress,
						CueBanner: "192.168.1.2",
						Text:      "",
						OnTextChanged: func() {
						},
					},
					Label{
						Text: "Port:",
					},
					NumberEdit{
						AssignTo:    &consolePort,
						Value:       float64(5201),
						ToolTipText: "1~65535",
						MaxValue:    65535,
						MinValue:    1,
						OnValueChanged: func() {
							// addLink.Port = int(consolePort.Value())
						},
					},
					Label{
						Text: "Parallel Streams:",
					},
					NumberEdit{
						AssignTo:    &consoleStreams,
						Value:       float64(5),
						ToolTipText: "1~100",
						MaxValue:    100,
						MinValue:    1,
						OnValueChanged: func() {
							// addLink.Port = int(consolePort.Value())
						},
					},
					Label{
						Text: "Packet Length:",
					},
					NumberEdit{
						AssignTo:    &consoleLength,
						Value:       float64(5),
						ToolTipText: "1~100",
						MaxValue:    100,
						MinValue:    1,
						OnValueChanged: func() {
							// addLink.Port = int(consolePort.Value())
						},
					},
					Label{
						Text: "Protocol:",
					},
					ComboBox{
						AssignTo:     &consoleProtocol,
						CurrentIndex: 0,
						Model:        []string{"tcp", "udp"},
						OnCurrentIndexChanged: func() {
							// addLink.Protocol = consoleProtocol.Text()
						},
					},
					Label{
						Text: "IP Version:",
					},
					ComboBox{
						AssignTo:     &consoleVersion,
						CurrentIndex: 0,
						Model:        []string{"auto", "ipv4", "ipv6"},
						OnCurrentIndexChanged: func() {
							// addLink.Protocol = consoleProtocol.Text()
						},
					},
					Label{
						Text: "Duration:",
					},
					NumberEdit{
						AssignTo:    &consoleTimeout,
						Value:       float64(30),
						ToolTipText: "1~120",
						MaxValue:    120,
						MinValue:    1,
						OnValueChanged: func() {
							// addLink.Port = int(consolePort.Value())
						},
					},
					Label{
						Text: "Bandwidth:",
					},
					LineEdit{
						Text:       "700 Mbit",
						ReadOnly:   true,
						Persistent: true,
					},
					Label{
						Text: "Statistics:",
					},
					LineEdit{
						Text:       "MAX:700 MIN:600 AVG:500",
						ReadOnly:   true,
						Persistent: true,
					},

					VSpacer{},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{

							PushButton{
								AssignTo: &buttonStart,
								Text:     "Start",
								OnClicked: func() {
									buttonStart.SetEnabled(false)
									buttonStop.SetEnabled(true)
								},
							},
							PushButton{
								AssignTo: &buttonStop,
								Text:     "Stop",
								Enabled:  false,
								OnClicked: func() {
									buttonStart.SetEnabled(true)
									buttonStop.SetEnabled(false)
								},
							},
						},
					},
				},
			},
		},
	}.Run()

	if err != nil {
		log.Fatalf(err.Error())
	} else {
		mainWin.Close()
	}
}
